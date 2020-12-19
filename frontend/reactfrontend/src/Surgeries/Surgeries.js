import React, { Component } from "react";
import dayjs from "dayjs";
import PatientSearchOrNew from "../PatientSearchOrNew/PatientSearchOrNew";
import { withTranslation } from "react-i18next";

class Surgeries extends Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedDoctorValue: "selectAnOption",
      selectedDoctor: null,
      selectedPatient: props.currentPatient,
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
      window.alert(this.props.t("FillAllFields"));
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
    };
    let ajax = new XMLHttpRequest();
    let url = "/api/surgeries";
    ajax.open("POST", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify(surgery));
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        window.alert(this.props.t("SurgerySaved"));
        this.clearForm();
      }
      if (ajax.readyState === 4 && ajax.status !== 200) {
        window.alert(this.props.t("CouldNotCompleteOperation"));
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
  setProposedSurgery(s) {
    this.setState({
      proposedSurgery: s,
    });
  }

  localSubmitHandler(e) {
    e.preventDefault();
  }
  render() {
    const { t } = this.props;
    return (
      <div>
        <h1>{t("Surgeries")}</h1>
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
              <label htmlFor="proposedSurgery">{t("ProposedSurgery")}: </label>
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
                  this.saveSurgery();
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
export default withTranslation()(Surgeries);
