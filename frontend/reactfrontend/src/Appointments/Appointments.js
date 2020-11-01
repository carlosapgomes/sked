import React, { Component } from "react";

export default class Appointments extends Component {
  constructor(props) {
    super(props);

    this.state = {
      selectedDoctor: null,
      searchField: null,
      patientSearchResult: [],
      selectedPatient: null,
    };
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
              <select name="doctor" id="doctor">
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
          <input type="date" id="apptmtdate" name="apptmtdate" />
          <label htmlFor="pctsearch">Patient: </label>
          <input
            type="text"
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
            <select name="searchresult" id="searchresult" size="4">
              {this.state.patientSearchResult.map((p) => {
                return (
                  <option key={p.id} value={p.id}>
                    {p.name}
                  </option>
                );
              })}
            </select>
          </div>
        </form>
      </div>
    );
  }
}
