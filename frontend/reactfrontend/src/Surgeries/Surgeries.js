import React, { Component } from "react";
import PatientSearch from "../PatientSearch/PatientSearch";
import dayjs from "dayjs";

export default class Surgeries extends Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedDoctorValue: "selectAnOption",
      selectedDoctor: null,
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
  setSelectedPatient(p) {
    if (!p) {
      this.setState({
        selectedPatient: null,
      });
    } else {
      this.setState({
        selectedPatient: { ...p },
      });
    }
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

  localSubmitHandler(e) {
    e.preventDefault();
  }
  render() {
    return (
      <div>
        <h1>Surgeries</h1>
        <section>
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
            <PatientSearch
              setSelectedPatient={(p) => {
                this.setSelectedPatient(p);
              }}
            />
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
        </section>
      </div>
    );
  }
}
