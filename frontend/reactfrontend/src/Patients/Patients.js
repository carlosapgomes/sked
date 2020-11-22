import React, { Component } from "react";
import PatientSearch from "../PatientSearch/PatientSearch";
import "./Patients.css";
export default class Patients extends Component {
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
        let data = JSON.parse(ajax.response);
        if (data) {
          this.props.updateCurrentPatient({ ...data });
          console.log(data);
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
  clearForm() {
    this.setState({
      id: "",
      name: "",
      address: "",
      city: "",
      state: "",
      phones: [],
      showUpdateButton: false,
    });
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
  setSelectedPatient(p) {
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
  render() {
    return (
      <div>
        <h1>Patients</h1>
        <section className="PatientsSection">
          <div className="ColumnItem">
            <PatientSearch
              setSelectedPatient={(p) => {
                this.setSelectedPatient(p);
              }}
            />
          </div>
          <br />
          <div className="ColumnItem">
            <form
              acceptCharset="utf-8"
              onSubmit={(e) => {
                this.localSubmitHandler(e);
              }}
            >
              <div>
                <label htmlFor="name">Name: </label>
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
                <label htmlFor="address">Address: </label>
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
                <label htmlFor="city">City: </label>
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
                <label htmlFor="state">State: </label>
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
                <label htmlFor="phones">Phones: </label>
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
                  Save
                </button>
                <button
                  hidden={!this.state.showUpdateButton}
                  onClick={() => {
                    this.updatePatient();
                  }}
                >
                  Update
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
        </section>
      </div>
    );
  }
}
