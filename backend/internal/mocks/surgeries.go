package mocks

import (
	"time"

	"carlosapgomes.com/sked/internal/surgery"
)

// mocks surgery services and repository

// SurgeryMockRepo is a mocked surgery repository
type SurgeryMockRepo struct {
	aDb []surgery.Surgery
}

//NewSurgeryRepo returns a mocked repository
func NewSurgeryRepo() *SurgeryMockRepo {
	var db []surgery.Surgery
	return &SurgeryMockRepo{
		db,
	}
}

// Create
func (r SurgeryMockRepo) Create(surgery surgery.Surgery) (*string, error) {
	var id string
	id = surgery.ID
	return &id, nil
}

// Update
func (r SurgeryMockRepo) Update(surgery surgery.Surgery) (*string, error) {
	var id string
	id = surgery.ID
	return &id, nil
}

// FindByID
func (r SurgeryMockRepo) FindByID(id string) (*surgery.Surgery, error) {
	surg := surgery.Surgery{
		ID:              "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		DateTime:        time.Now(),
		PatientName:     "John Doe",
		PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:      "Dr House",
		DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:           "some notes",
		ProposedSurgery: "saphenectomy",
		CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:       time.Now(),
	}
	if id == surg.ID {
		return &surg, nil
	} else {
		return nil, surgery.ErrNoRecord
	}
}

// FindByPatientID
func (r SurgeryMockRepo) FindByPatientID(patientID string) ([]surgery.Surgery, error) {
	surg := surgery.Surgery{
		ID:              "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		DateTime:        time.Now(),
		PatientName:     "John Doe",
		PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:      "Dr House",
		DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:           "some notes",
		ProposedSurgery: "saphenectomy",
		CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:       time.Now(),
	}
	if patientID == surg.PatientID {
		surgs := []surgery.Surgery{
			surg,
		}
		return surgs, nil
	} else {
		return nil, surgery.ErrNoRecord
	}
}

// FindByDoctorID
func (r SurgeryMockRepo) FindByDoctorID(doctorID string) ([]surgery.Surgery, error) {
	surg := surgery.Surgery{
		ID:          "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		DateTime:    time.Now(),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Now(),
	}
	if doctorID == surg.DoctorID {
		surgs := []surgery.Surgery{
			surg,
		}
		return surgs, nil
	} else {
		return nil, surgery.ErrNoRecord
	}
}

// FindByDate
func (r SurgeryMockRepo) FindByDate(dateTime time.Time) ([]surgery.Surgery, error) {
	surg := surgery.Surgery{
		ID:          "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		DateTime:    time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
		PatientName: "John Doe",
		PatientID:   "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:  "Dr House",
		DoctorID:    "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:       "some notes",
		CreatedBy:   "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:   time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}
	searchY, searchM, searchD := dateTime.Date()
	appointmtY, appointmtM, appointmtD := surg.DateTime.Date()
	if (searchY == appointmtY) && (searchM == appointmtM) && (searchD == appointmtD) {
		appointmts := []surgery.Surgery{
			surg,
		}
		return appointmts, nil
	} else {
		return nil, surgery.ErrNoRecord
	}
}

// GetAll
func (r SurgeryMockRepo) GetAll(cursor string, after bool, pgSize int) ([]surgery.Surgery, bool, error) {
	var db []surgery.Surgery
	db = append(db, surgery.Surgery{
		ID:              "e521798b-9f33-4a10-8b2a-9677ed1cd1ae",
		DateTime:        time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
		PatientName:     "John Doe",
		PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:      "Dr House",
		DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:           "some notes",
		ProposedSurgery: "saphenectomy",
		CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, surgery.Surgery{
		ID:              "5e6f7cd1-d8d2-40cd-97a3-aca01a93bfde",
		DateTime:        time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC),
		PatientName:     "John Doe",
		PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:      "Dr House",
		DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:           "some notes",
		ProposedSurgery: "saphenectomy",
		CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, surgery.Surgery{
		ID:              "7fef3c47-a01a-42a6-ac45-27a440596751",
		DateTime:        time.Date(2020, 9, 8, 12, 0, 0, 0, time.UTC),
		PatientName:     "John Doe",
		PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:      "Dr House",
		DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:           "some notes",
		ProposedSurgery: "saphenectomy",
		CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, surgery.Surgery{
		ID:              "19f66dc6-b5c8-497b-bba2-b982bb85ded8",
		DateTime:        time.Date(2020, 9, 9, 12, 0, 0, 0, time.UTC),
		PatientName:     "John Doe",
		PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:      "Dr House",
		DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:           "some notes",
		ProposedSurgery: "saphenectomy",
		CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:       time.Date(2020, 9, 10, 12, 0, 0, 0, time.UTC),
	}, surgery.Surgery{
		ID:              "723e2fa0-70a9-4c20-89d9-b5f69405b772",
		DateTime:        time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
		PatientName:     "John Doe",
		PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:      "Dr House",
		DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:           "some notes",
		ProposedSurgery: "saphenectomy",
		CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	}, surgery.Surgery{
		ID:              "5792340a-8c35-4183-a388-2459a8e0295a",
		DateTime:        time.Date(2020, 9, 11, 12, 0, 0, 0, time.UTC),
		PatientName:     "John Doe",
		PatientID:       "22070f56-5d52-43f0-9f59-5de61c1db506",
		DoctorName:      "Dr House",
		DoctorID:        "f06244b9-97e5-4f1a-bae0-3b6da7a0b604",
		Notes:           "some notes",
		ProposedSurgery: "saphenectomy",
		CreatedBy:       "10b9ad06-e86d-4a85-acb1-d7e268d1f21a",
		CreatedAt:       time.Date(2020, 9, 6, 12, 0, 0, 0, time.UTC),
	},
	)
	var res []surgery.Surgery
	var hasMore bool
	hasMore = false
	var respSize int
	if cursor == "" {
		if pgSize >= len(db) {
			respSize = len(db)
			hasMore = false
		} else {
			respSize = pgSize
			hasMore = true
		}
		for i := 0; i < respSize; i++ {
			res = append(res, db[i])
		}
		return res, hasMore, nil
	}
	pos := r.findPos(db, cursor)
	if pos == -1 {
		return nil, false, surgery.ErrNoRecord
	}

	if after {
		start := pos + 1
		for i := start; i < (start + pgSize); i++ {
			res = append(res, db[i])
		}
		if (len(db) - pos) > pgSize {
			hasMore = true
		}
	} else {
		start := pos - pgSize
		if start < 0 {
			start = 0
		}
		for i := start; i <= (pos - 1); i++ {
			res = append(res, db[i])
		}
		if pos > pgSize {
			hasMore = true
		}
	}
	return res, hasMore, nil
}
func (r SurgeryMockRepo) findPos(surgs []surgery.Surgery, id string) int {
	for i, el := range surgs {
		if el.ID == id {
			return i
		}
	}
	return -1
}
