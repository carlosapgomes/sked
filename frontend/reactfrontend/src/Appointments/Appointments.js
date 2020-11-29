import React, { Component } from "react";
import { withTranslation } from "react-i18next";
import PatientSearchOrNew from "../PatientSearchOrNew/PatientSearchOrNew";
import dayjs from "dayjs";

class Appointments extends Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedDoctorValue: "selectAnOption",
      selectedDoctor: null,
      selectedPatient: props.currentPatient,
      selectedDate: "",
      selectedTime: "08:00",
      notes: "",
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
    });
  }
  clearForm() {
    this.resetInitialState();
  }
  saveAppointment() {
    if (
      this.state.selectedDoctor == null ||
      this.state.selectedPatient === null ||
      this.state.selectedDate === "" ||
      this.state.selectedTime === ""
    ) {
      window.alert("Fill all fields");
      return;
    }
    let dateTime = dayjs(
      this.state.selectedDate + " " + this.state.selectedTime
    );
    let appointmt = {
      dateTime: dateTime.format(),
      patientName: this.state.selectedPatient.name,
      patientID: this.state.selectedPatient.id,
      doctorName: this.state.selectedDoctor.name,
      doctorID: this.state.selectedDoctor.id,
      notes: this.state.notes,
    };
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/appointments";
    ajax.open("POST", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify(appointmt));
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        window.alert("Appointment saved");
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
  updateCurrentPatient(p) {
    if (!p) {
      this.setState({
        selectedPatient: null,
      });
      this.props.updateCurrentPatient(null);
    } else {
      this.setState({
        selectedPatient: { ...p },
      });
      this.props.updateCurrentPatient({ ...p });
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
  updateSearchField(s) {
    this.setState({
      searchField: s,
    });
  }

  localSubmitHandler(e) {
    e.preventDefault();
  }
  render() {
    const { t } = this.props;
    return (
      <div>
        <h1>{t("Appointments")}</h1>
        <section>
          <PatientSearchOrNew
            currentPatient={this.state.currentPatient}
            updateCurrentPatient={(p) => {
              this.updateCurrentPatient(p);
            }}
          />
        </section>
        <section>
          <form
            onSubmit={(e) => {
              this.localSubmitHandler(e);
            }}
          >
            <label>{t("Doctor")}: </label>
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
                    -- {t("SelectAnOption")} --{" "}
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
            <label htmlFor="apptmtdatetime">{t("DateTime")}: </label>
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
            <div>
              <label htmlFor="notes">{t("Notes")}: </label>
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
                  this.saveAppointment();
                }}
              >
                {t("Save")}
              </button>
              &nbsp;&nbsp;
              <button
                onClick={() => {
                  this.clearForm();
                }}
              >
                {t("Clear")}
              </button>
            </div>
          </form>
        </section>
      </div>
    );
  }
}
export default withTranslation()(Appointments);
