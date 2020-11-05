import React, { Component } from "react";

export default class Patients extends Component {
  constructor(props) {
    super(props);

    this.state = {
      name: "",
      address: "",
      city: "",
      state: "",
      phones: [],
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
    console.log(patient);
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/patients";
    ajax.open("POST", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify(patient));
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
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
      name: "",
      address: "",
      city: "",
      state: "",
      phones: [],
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
  localSubmitHandler(e) {
    e.preventDefault();
  }
  render() {
    return (
      <div>
        <h1>New Patient</h1>
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
              onClick={() => {
                this.savePatient();
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
