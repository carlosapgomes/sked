import React, { Component } from "react";
import { withTranslation } from "react-i18next";
import "./PatientSearchOrNew.css";
class PatientSearchOrNew extends Component {
  constructor(props) {
    super(props);

    this.state = {
      id: "",
      name: "",
      address: "",
      city: "",
      state: "",
      phones: [],
      showUpdateButton: false,
      searchField: "",
      patientSearchResult: [],
      selectedPatientValue: "selectAnOption",
      selectedPatient: null,
      showNewPatientForm: false,
    };
  }
  savePatient() {
    if (
      this.state.name === "" ||
      this.state.city === "" ||
      this.state.phones.length === 0
    ) {
      window.alert("Please, fill in at least name, city and phones");
      return;
    }
    let patient = {
      name: this.state.name,
      address: this.state.address,
      city: this.state.city,
      state: this.state.state,
      phones: [...this.state.phones],
    };
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/patients";
    ajax.open("POST", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify(patient));
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        if (ajax.response) {
          let data = JSON.parse(ajax.response);
          this.props.updateCurrentPatient({ ...data });
        }
        window.alert("Patient saved");
        this.clearForm();
      }
      if (ajax.readyState === 4 && ajax.status !== 200) {
        window.alert("Could not complete operation");
        console.log(ajax.responseText);
      }
    };
  }
  updatePatient() {
    if (
      this.state.id === "" ||
      this.state.name === "" ||
      this.state.city === "" ||
      this.state.phones.length === 0
    ) {
      window.alert("Please, fill in at least name, city and phones");
      return;
    }
    let patient = {
      id: this.state.id,
      name: this.state.name,
      address: this.state.address,
      city: this.state.city,
      state: this.state.state,
      phones: [...this.state.phones],
    };
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/patients/" + patient.id;
    ajax.open("PUT", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify(patient));
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        if (ajax.response) {
          let data = JSON.parse(ajax.response);
          this.props.updateCurrentPatient({ ...data });
        }
        window.alert("Patient updated");
        this.clearForm();
      }
      if (ajax.readyState === 4 && ajax.status !== 200) {
        window.alert("Could not complete operation");
        console.log(ajax.responseText);
      }
    };
  }
  clearForm() {
    this.setState({
      id: "",
      name: "",
      address: "",
      city: "",
      state: "",
      phones: [],
      showUpdateButton: false,
      searchField: "",
      patientSearchResult: [],
      selectedPatientValue: "selectAnOption",
      selectedPatient: null,
      showNewPatientForm: false,
    });
  }
  setSelectedPatient(e) {
    let idx = e.target.selectedIndex - 1;
    this.setState({
      selectedPatientValue: this.state.patientSearchResult[idx].id,
      selectedPatient: { ...this.state.patientSearchResult[idx] },
      searchField: "",
    });
    this.setState({
      id: this.state.patientSearchResult[idx].id,
      name: this.state.patientSearchResult[idx].name,
      address: this.state.patientSearchResult[idx].address,
      city: this.state.patientSearchResult[idx].city,
      state: this.state.patientSearchResult[idx].state,
      phones: [...this.state.patientSearchResult[idx].phones],
      showUpdateButton: true,
    });
    // update selectedPatient on parent component
    //this.props.setSelectedPatient({ ...this.state.patientSearchResult[idx] });
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
  setName(s) {
    this.setState({
      name: s,
    });
  }
  setAddress(s) {
    this.setState({
      address: s,
    });
  }
  setCity(s) {
    this.setState({
      city: s,
    });
  }
  setSt(s) {
    this.setState({
      state: s,
    });
  }
  setPhones(s) {
    let phones = [];
    phones.push(s);
    this.setState({
      phones: [...phones],
    });
  }
  OrigSetSelectedPatient(p) {
    if (!p) {
      this.clearForm();
    } else {
      this.setState({
        id: p.id,
        name: p.name,
        address: p.address,
        city: p.city,
        state: p.state,
        phones: [...p.phones],
        showUpdateButton: true,
      });
    }
  }
  localSubmitHandler(e) {
    e.preventDefault();
  }
  toggleNewPatientForm() {
    this.setState({
      showNewPatientForm: !this.state.showNewPatientForm,
    });
  }
  render() {
    const { t } = this.props;
    return (
      <div>
        <label>{t("Patient")}:</label>
        <section className="PatientsSection">
          <br />
          <div className="ColumnItem">
            <form
              acceptCharset="utf-8"
              onSubmit={(e) => {
                this.localSubmitHandler(e);
              }}
            >
              <div>
                <input
                  type="text"
                  placeholder={t("Name")}
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
                &nbsp;&nbsp;
                <button
                  onClick={() => {
                    this.toggleNewPatientForm();
                  }}
                >
                  {t("New")}
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
                      -- {t("SelectAnOption")} --{" "}
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
              </div>{" "}
              <div
                style={{
                  display: this.state.showNewPatientForm ? "block" : "none",
                }}
              >
                <div>
                  <label htmlFor="name">{t("Name")}: </label>
                  <input
                    type="text"
                    name="name"
                    id="name"
                    value={this.state.name}
                    onChange={(e) => {
                      this.setName(e.target.value);
                    }}
                  />{" "}
                </div>
                <div>
                  <label htmlFor="address">{t("Address")}: </label>
                  <input
                    type="text"
                    name="address"
                    id="address"
                    value={this.state.address}
                    onChange={(e) => {
                      this.setAddress(e.target.value);
                    }}
                  />
                </div>
                <div>
                  <label htmlFor="city">{t("City")}: </label>
                  <input
                    type="text"
                    name="city"
                    id="city"
                    value={this.state.city}
                    onChange={(e) => {
                      this.setCity(e.target.value);
                    }}
                  />
                </div>
                <div>
                  <label htmlFor="state">{t("State")}: </label>
                  <input
                    type="text"
                    name="state"
                    id="state"
                    value={this.state.state}
                    onChange={(e) => {
                      this.setSt(e.target.value);
                    }}
                  />
                </div>
                <div>
                  <label htmlFor="phones">{t("Phones")}: </label>
                  <input
                    type="tel"
                    name="phones"
                    id="phones"
                    value={this.state.phones.toString()}
                    onChange={(e) => {
                      this.setPhones(e.target.value);
                    }}
                  />
                </div>
                <div>
                  <button
                    hidden={this.state.showUpdateButton}
                    onClick={() => {
                      this.savePatient();
                    }}
                  >
                    {t("Save")}
                  </button>
                  <button
                    hidden={!this.state.showUpdateButton}
                    onClick={() => {
                      this.updatePatient();
                    }}
                  >
                    {t("Update")}
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
              </div>
            </form>
          </div>
        </section>
      </div>
    );
  }
}

export default withTranslation()(PatientSearchOrNew);
