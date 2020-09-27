package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"carlosapgomes.com/sked/internal/appointment"
	"carlosapgomes.com/sked/internal/surgery"
	"carlosapgomes.com/sked/internal/user"
	uuid "github.com/satori/go.uuid"
)

// surgeryService provides implementation of surgery domain interface
type surgeryService struct {
	repo    surgery.Repository
	userSvc user.Service
}

// NewSurgeryService returns a surgery Service instance
func NewSurgeryService(repo surgery.Repository, userSvc user.Service) surgery.Service {
	return &surgeryService{
		repo,
		userSvc,
	}
}

// Create - creates a new surgery and returns its uuid
func (s *surgeryService) Create(dateTime time.Time, patientName, patientID, doctorName, doctorID, notes, proposedSurgery, createdBy string) (*string, error) {
	// validate ID format (uuidV4)
	ptID, err := uuid.FromString(patientID)
	if err != nil {
		return nil, surgery.ErrInvalidInputSyntax
	}
	docID, err := uuid.FromString(doctorID)
	if err != nil {
		return nil, surgery.ErrInvalidInputSyntax
	}
	createdByID, err := uuid.FromString(createdBy)
	if err != nil {
		return nil, surgery.ErrInvalidInputSyntax
	}
	// user with ID == doctorID must have RoleDoctor
	userDoc, err := s.userSvc.FindByID(doctorID)
	if err != nil {
		return nil, err
	}
	isDoc := false
	for i := range userDoc.Roles {
		if userDoc.Roles[i] == user.RoleDoctor {
			isDoc = true
		}
	}
	if !isDoc {
		return nil, appointment.ErrInvalidInputSyntax
	}
	// if appointment is not created by the same doctor,
	// it can only be created by a clerk or admin
	if doctorID != createdBy {
		createdByUser, err := s.userSvc.FindByID(createdBy)
		if err != nil {
			return nil, err
		}
		isClerkOrAdmin := false
		for i := range createdByUser.Roles {
			if (createdByUser.Roles[i] == user.RoleAdmin) ||
				(createdByUser.Roles[i] == user.RoleClerk) {
				isClerkOrAdmin = true
			}
			if !isClerkOrAdmin {
				return nil, appointment.ErrInvalidInputSyntax
			}
		}
	}

	uid := uuid.NewV4()
	dt := dateTime.UTC()
	newSurgery := surgery.Surgery{
		ID:              uid.String(),
		DateTime:        dt,
		PatientName:     patientName,
		PatientID:       ptID.String(),
		DoctorName:      doctorName,
		DoctorID:        docID.String(),
		Notes:           notes,
		ProposedSurgery: proposedSurgery,
		Canceled:        false,
		Done:            false,
		CreatedBy:       createdByID.String(),
		CreatedAt:       time.Now().UTC(),
	}

	id, err := s.repo.Create(newSurgery)
	if err != nil {
		return nil, err
	}
	if (id != nil) && (*id != newSurgery.ID) {
		return nil, errors.New("New surgery creation: returned repository ID not equal to new surgery ID")
	}
	return id, err
}

// Update - updates a surgery
func (s *surgeryService) Update(surg surgery.Surgery) (*string, error) {
	// get original surgery
	original, err := s.repo.FindByID(surg.ID)
	if err != nil {
		fmt.Print("could not find surgery\n")
		return nil, surgery.ErrNoRecord
	}
	_, err = uuid.FromString(surg.UpdatedBy)
	if err != nil {
		return nil, surgery.ErrInvalidInputSyntax
	}
	original.Notes = surg.Notes
	original.ProposedSurgery = surg.ProposedSurgery
	original.DateTime = surg.DateTime
	original.Canceled = surg.Canceled
	original.Done = surg.Done
	original.UpdatedAt = time.Now().UTC()
	original.UpdatedBy = surg.UpdatedBy

	id, err := s.repo.Update(*original)
	if err != nil {
		return nil, err
	}
	if (id != nil) && (*id != original.ID) {
		return nil, errors.New("Surgery update: returned repository ID not equal to new surgery ID")
	}
	return id, nil
}

// FindByID - look for a surgery by its uuid
func (s *surgeryService) FindByID(id string) (*surgery.Surgery, error) {
	_, err := uuid.FromString(id)
	if err != nil {
		return nil, surgery.ErrInvalidInputSyntax
	}
	surgery, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return surgery, nil
}

// FindFindByPatientID - look for surgeries by its patientID
func (s *surgeryService) FindByPatientID(patientID string) ([]*surgery.Surgery, error) {
	_, err := uuid.FromString(patientID)
	if err != nil {
		return nil, surgery.ErrInvalidInputSyntax
	}
	surgeries, err := s.repo.FindByPatientID(patientID)
	if err != nil {
		return nil, err
	}
	return surgeries, nil
}

// FindByDoctorID - look for surgeries by doctorID
func (s *surgeryService) FindByDoctorID(doctorID string) ([]*surgery.Surgery, error) {
	_, err := uuid.FromString(doctorID)
	if err != nil {
		return nil, surgery.ErrInvalidInputSyntax
	}
	surgeries, err := s.repo.FindByDoctorID(doctorID)
	if err != nil {
		return nil, err
	}
	return surgeries, nil
}

// FindByDate - look for surgeries by date
func (s *surgeryService) FindByDate(dateTime time.Time) ([]*surgery.Surgery, error) {
	surgeries, err := s.repo.FindByDate(dateTime)
	if err != nil {
		return nil, err
	}
	return surgeries, nil
}

// GetAll - return all surgeries
func (s *surgeryService) GetAll(before string, after string, pgSize int) (*surgery.Cursor, error) {
	var surgeriesResp surgery.Cursor
	var err error
	var aList *[]surgery.Surgery
	if pgSize <= 0 {
		pgSize = 15
	}

	switch {
	case (before != "" && after != ""):
		// if both (before & after) are present, returns error
		return nil, surgery.ErrInvalidInputSyntax
	case (before == "" && after == ""):
		// if they are empty
		// get default list and page size
		aList, surgeriesResp.HasBefore, err = s.repo.
			GetAll("", false, pgSize)
		if err != nil {
			return nil, err
		}
		if aList != nil {
			for _, a := range *aList {
				surgeriesResp.Surgeries = append(surgeriesResp.Surgeries, a)
			}
		}
		if len(surgeriesResp.Surgeries) > 0 {
			surgeriesResp.Before = base64.StdEncoding.
				EncodeToString([]byte(surgeriesResp.Surgeries[len(surgeriesResp.Surgeries)-1].ID))
		} else {
			surgeriesResp.Before = ""
		}
		surgeriesResp.After = ""
		surgeriesResp.HasAfter = false
		// and return values
	case (before != ""):
		// if before is present,
		// get a before list
		c, err := base64.StdEncoding.DecodeString(before)
		if err != nil {
			return nil, err
		}
		cursor := string(c)
		aList, surgeriesResp.HasBefore, err = s.repo.GetAll(cursor, false, pgSize)
		if err != nil {
			return nil, err
		}
		if aList != nil {
			for _, a := range *aList {
				surgeriesResp.Surgeries = append(surgeriesResp.Surgeries, a)
			}
		}
		if len(surgeriesResp.Surgeries) > 0 {
			befCursor := base64.StdEncoding.EncodeToString([]byte(surgeriesResp.Surgeries[len(surgeriesResp.Surgeries)-1].ID))
			surgeriesResp.Before = befCursor
		} else {
			surgeriesResp.Before = ""
		}
		// test for 'after data' from the requested cursor
		// fill the response fields
		_, surgeriesResp.HasAfter, err = s.repo.GetAll(cursor, true, pgSize)
		if surgeriesResp.HasAfter {
			surgeriesResp.After = base64.StdEncoding.EncodeToString([]byte(before))
		} else {
			surgeriesResp.After = ""
		}
		// and return it
	case (after != ""):
		// if after is present,
		// get an after list
		c, err := base64.StdEncoding.DecodeString(after)
		if err != nil {
			return nil, err
		}
		cursor := string(c)
		aList, surgeriesResp.HasAfter, err = s.repo.
			GetAll(cursor, true, pgSize)
		// and return it
		if aList != nil {
			for _, a := range *aList {
				surgeriesResp.Surgeries = append(surgeriesResp.Surgeries, a)
			}
		}
		if len(surgeriesResp.Surgeries) > 0 {
			surgeriesResp.After = base64.StdEncoding.EncodeToString([]byte(surgeriesResp.Surgeries[0].ID))
		}
		// test for 'before data' from the requested cursor
		// fill the response fields
		_, surgeriesResp.HasBefore, err = s.repo.
			GetAll(cursor, false, pgSize)
		if surgeriesResp.HasBefore {
			surgeriesResp.Before = base64.StdEncoding.EncodeToString([]byte(after))
		}
	}
	return &surgeriesResp, nil
}
