import React, { Component } from "react";

export default class Appointments extends Component {
  constructor(props) {
    super(props);
    this.state = {
      initialSelectedValue: "selectAnOption",
      selectedDoctor: null,
      searchField: "",
      patientSearchResult: [],
      selectedPatient: null,
      selectedDate: "",
    };
  }
  resetInitialState() {
    this.setState({
      initialSelectedValue: "selectAnOption",
      selectedDoctor: null,
      searchField: "",
      patientSearchResult: [],
      selectedPatient: null,
      selectedDate: "",
    });
  }
  clearForm() {
    this.resetInitialState();
  }
  saveAppointment() {}
  setSelectedPatient(e) {
    let idx = e.target.selectedIndex;
    this.setState({
      selectedPatient: { ...this.state.patientSearchResult[idx] },
    });
  }
  setDate(e) {
    this.setState({
      selectedDate: e.target.value,
    });
  }
  setSelectedDoctor(e) {
    let idx = e.target.selectedIndex - 1;
    if (idx >= 0) {
      this.setState({
        selectedDoctor: { ...this.props.doctors[idx] },
      });
    }
  }
  updateSearchField(s) {
    this.setState({
      searchField: s,
    });
  }
  searchPatient(e) {
    e.preventDefault();
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
        var data = JSON.parse(ajax.responseText);
        if (typeof data != "undefined") {
          this.setState({
            patientSearchResult: data,
          });
        }
      }
    };
  }

  localSubmitHandler(e) {
    e.preventDefault();
  }
  render() {
    return (
      <div>
        <h1>Appointments</h1>
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
                value={this.state.initialSelectedValue}
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
          <label htmlFor="apptmtdate">Data: </label>
          <input
            type="date"
            id="apptmtdate"
            name="apptmtdate"
            value={this.state.selectedDate}
            onChange={(e) => {
              this.setDate(e);
            }}
          />
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
            onClick={(e) => {
              this.searchPatient(e);
            }}
          >
            Search
          </button>
          <div hidden={this.state.patientSearchResult.length <= 0}>
            <select
              name="searchresult"
              id="searchresult"
              size="4"
              onChange={(e) => {
                this.setSelectedPatient(e);
              }}
            >
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
            <button
              onClick={() => {
                this.saveAppointment();
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
