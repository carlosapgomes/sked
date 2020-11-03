import React, { Component } from "react";
import dayjs from "dayjs";

export default class Surgeries extends Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedDoctorValue: "selectAnOption",
      selectedDoctor: null,
      searchField: "",
      patientSearchResult: [],
      selectedPatientValue: "selectAnOption",
      selectedPatient: null,
      selectedDate: "",
      selectedTime: "08:00",
      notes: "",
      proposedSurgery: "",
    };
  }
  resetInitialState() {
    this.setState({
      selectedDoctorValue: "selectAnOption",
      selectedDoctor: null,
      searchField: "",
      patientSearchResult: [],
      selectedPatientValue: "selectAnOption",
      selectedPatient: null,
      selectedDate: "",
      selectedTime: "08:00",
      notes: "",
      proposedSurgery: "",
    });
  }
  clearForm() {
    this.resetInitialState();
  }
  saveSurgery() {
    if (
      this.state.selectedDoctor == null ||
      this.state.selectedPatient === null ||
      this.state.selectedDate === "" ||
      this.state.selectedTime === "" ||
      this.state.proposedSurgery === ""
    ) {
      window.alert("Please, fill all fields");
      return;
    }
    let dateTime = dayjs(
      this.state.selectedDate + " " + this.state.selectedTime
    );
    let surgery = {
      dateTime: dateTime.format(),
      patientName: this.state.selectedPatient.name,
      patientID: this.state.selectedPatient.id,
      doctorName: this.state.selectedDoctor.name,
      doctorID: this.state.selectedDoctor.id,
      notes: this.state.notes,
      proposedSurgery: this.state.proposedSurgery,
      createdBy: this.props.currentUser.uid,
    };
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/surgeries";
    ajax.open("POST", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify(surgery));
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        window.alert("Surgery saved");
        this.clearForm();
      }
      if (ajax.readyState === 4 && ajax.status !== 200) {
        window.alert("Could not complete operation");
        console.log(ajax.responseText);
      }
    };
  }
  setSelectedDoctor(e) {
    let idx = e.target.selectedIndex - 1;
    if (idx >= 0) {
      this.setState({
        selectedDoctorValue: this.props.doctors[idx].id,
        selectedDoctor: { ...this.props.doctors[idx] },
      });
    }
  }
  setSelectedPatient(e) {
    let idx = e.target.selectedIndex - 1;
    this.setState({
      selectedPatientValue: this.state.patientSearchResult[idx].id,
      selectedPatient: { ...this.state.patientSearchResult[idx] },
    });
  }
  setTime(e) {
    this.setState({
      selectedTime: e.target.value,
    });
  }
  setDate(e) {
    this.setState({
      selectedDate: e.target.value,
    });
  }
  setNotes(s) {
    this.setState({
      notes: s,
    });
  }
  setProposedSurgery(s) {
    this.setState({
      proposedSurgery: s,
    });
  }
  updateSearchField(s) {
    this.setState({
      searchField: s,
    });
  }
  searchPatient() {
    if (!this.state.searchField || this.state.searchField.length < 3) {
      return;
    }
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/patients?name=" + this.state.searchField;
    ajax.open("GET", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send();
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        if (!ajax.responseText) {
          window.alert("Could not find any patient");
        } else {
          let data = JSON.parse(ajax.responseText);
          if (data) {
            this.setState({
              patientSearchResult: [...data],
            });
          }
        }
      }
      if (ajax.readyState === 4 && ajax.status !== 200) {
        window.alert("Could not complete the operation");
        console.log(ajax.responseText);
      }
    };
  }

  localSubmitHandler(e) {
    e.preventDefault();
  }
  render() {
    return (
      <div>
        <h1>Surgeries</h1>
        <form
          onSubmit={(e) => {
            this.localSubmitHandler(e);
          }}
        >
          <label>Doctor: </label>
          {this.props.doctors.length > 0 ? (
            <div>
              <select
                name="doctor"
                id="doctor"
                onChange={(e) => {
                  this.setSelectedDoctor(e);
                }}
                value={this.state.selectedDoctorValue}
              >
                <option
                  hidden
                  disabled
                  defaultValue
                  value="selectAnOption"
                  style={{ display: "none" }}
                >
                  {" "}
                  -- select an option --{" "}
                </option>
                {this.props.doctors.map((d) => {
                  return (
                    <option key={d.id} value={d.id}>
                      {d.name}
                    </option>
                  );
                })}
              </select>
            </div>
          ) : null}
          <label htmlFor="apptmtdatetime">Date/Time: </label>
          <div id="apptmtdatetime">
            <input
              type="date"
              id="apptmtdate"
              name="apptmtdate"
              value={this.state.selectedDate}
              onChange={(e) => {
                this.setDate(e);
              }}
            />
            <input
              type="time"
              value={this.state.selectedTime}
              name="time"
              id="time"
              onChange={(e) => {
                this.setTime(e);
              }}
            />
          </div>
          <label htmlFor="pctsearch">Patient: </label>
          <input
            type="text"
            value={this.state.searchField}
            onChange={(e) => {
              this.updateSearchField(e.target.value);
            }}
            id="pctsearch"
            name="pctsearch"
          />
          <button
            onClick={() => {
              this.searchPatient();
            }}
          >
            Search
          </button>
          <div hidden={this.state.patientSearchResult.length <= 0}>
            <select
              name="searchresult"
              id="searchresult"
              onChange={(e) => {
                this.setSelectedPatient(e);
              }}
              value={this.state.selectedPatientValue}
            >
              <option
                hidden
                disabled
                defaultValue
                value="selectAnOption"
                style={{ display: "none" }}
              >
                {" "}
                -- select an option --{" "}
              </option>
              {this.state.patientSearchResult.map((p) => {
                return (
                  <option key={p.id} value={p.id}>
                    {p.name}
                  </option>
                );
              })}
            </select>
          </div>
          <div>
            <label htmlFor="proposedSurgery">Proposed surgery: </label>
            <input
              type="text"
              value={this.state.proposedSurgery}
              name="proposedSurgery"
              id="proposedSurgery"
              onChange={(e) => {
                this.setProposedSurgery(e.target.value);
              }}
            />
          </div>
          <div>
            <label htmlFor="notes">Notes: </label>
            <textarea
              name="notes"
              id="notes"
              onChange={(e) => {
                this.setNotes(e.target.value);
              }}
              value={this.state.notes}
            ></textarea>
          </div>
          <div>
            <button
              onClick={() => {
                this.saveSurgery();
              }}
            >
              Save
            </button>
            <button
              onClick={() => {
                this.clearForm();
              }}
            >
              Clear
            </button>
          </div>
        </form>
      </div>
    );
  }
}
