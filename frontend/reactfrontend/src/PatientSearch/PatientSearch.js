import React, { Component } from "react";
import { withTranslation } from "react-i18next";

class PatientSearch extends Component {
  constructor(props) {
    super(props);

    this.state = {
      searchField: "",
      patientSearchResult: [],
      selectedPatientValue: "selectAnOption",
      selectedPatient: null,
    };
  }
  resetInitialState() {
    this.setState({
      searchField: "",
      patientSearchResult: [],
      selectedPatientValue: "selectAnOption",
      selectedPatient: null,
    });
  }
  setSelectedPatient(e) {
    let idx = e.target.selectedIndex - 1;
    this.setState({
      selectedPatientValue: this.state.patientSearchResult[idx].id,
      selectedPatient: { ...this.state.patientSearchResult[idx] },
      searchField: "",
    });
    // update selectedPatient on parent component
    this.props.setSelectedPatient({ ...this.state.patientSearchResult[idx] });
  }
  setSearchField(s) {
    this.setState({
      searchField: s,
    });
  }
  searchPatient() {
    if (!this.state.searchField || this.state.searchField.length < 3) {
      return;
    }
    let str = this.state.searchField.trim().split(/\s+/).join(" ");
    if (str.length < 3) {
      return;
    }
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/patients?name=" + str;
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
              selectedPatientValue: "selectAnOption",
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
  render() {
    const { t } = this.props;
    return (
      <div>
        <label htmlFor="pctsearch">{t("Patient")}: </label>
        <input
          type="text"
          value={this.state.searchField}
          onChange={(e) => {
            this.setSearchField(e.target.value);
          }}
          id="pctsearch"
          name="pctsearch"
        />
        <button
          onClick={() => {
            this.searchPatient();
          }}
        >
          {t("Search")}
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
      </div>
    );
  }
}
export default withTranslation()(PatientSearch);
