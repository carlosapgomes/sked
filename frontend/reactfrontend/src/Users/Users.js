import React, { Component } from "react";
import { withTranslation } from "react-i18next";

class Users extends Component {
  constructor(props) {
    super(props);

    this.state = {
      name: "",
      email: "",
      phone: "",
      roles: [],
    };
  }
  setUsername(s) {
    this.setState({
      name: s,
    });
  }
  setEmail(s) {
    this.setState({
      email: s,
    });
  }
  setPhone(s) {
    this.setState({
      phone: s,
    });
  }

  saveUser() {
    if (
      this.state.name === "" ||
      this.state.email === "" ||
      this.state.phone === ""
    ) {
      window.alert("Please, fill the requested data");
      return;
    }
    let newUser = {
      Name: this.state.name,
      Email: this.state.email,
      Phone: this.state.phone,
    };
    console.log(newUser);
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/users";
    ajax.open("POST", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify(newUser));
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        window.alert("User saved");
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
      email: "",
      phone: "",
    });
  }
  localSubmitHandler(e) {
    e.preventDefault();
  }
  updateRoles(e) {
    console.log(e.target);
  }
  render() {
    const { t } = this.props;
    return (
      <div>
        <h1>{t("NewUser")}</h1>
        <section>
          <form
            acceptCharset="utf-8"
            onSubmit={(e) => {
              this.localSubmitHandler(e);
            }}
          >
            <div>
              <label htmlFor="name">{t("Name")}: </label>
              <input
                type="text"
                value={this.state.name}
                name="name"
                id="name"
                onChange={(e) => {
                  this.setUsername(e.target.value);
                }}
              />
            </div>
            <div>
              <label htmlFor="email">{t("Email")}: </label>
              <input
                type="email"
                value={this.state.email}
                name="email"
                id="email"
                onChange={(e) => {
                  this.setEmail(e.target.value);
                }}
              />
            </div>
            <div>
              <label htmlFor="phone">{t("Phone")}: </label>
              <input
                type="tel"
                value={this.state.phone}
                name="phone"
                id="phone"
                onChange={(e) => {
                  this.setPhone(e.target.value);
                }}
              />
            </div>
            <div>
              <fieldset id="" class="">
                <legend>{t("Roles")}:</legend>
                <div>
                  <label htmlFor="clerk">{t("Clerk")}</label>
                  <input
                    type="checkbox"
                    value="Clerk"
                    name="clerk"
                    onChange={(e) => {
                      this.updateRoles(e);
                    }}
                  />
                </div>
                <div>
                  <label htmlFor="clerk">{t("Doctor")}</label>
                  <input
                    type="checkbox"
                    value="Doctor"
                    name="doctor"
                    onChange={(e) => {
                      this.updateRoles(e);
                    }}
                  />
                </div>
                <div>
                  <label htmlFor="clerk">{t("Admin")}</label>
                  <input
                    type="checkbox"
                    value="Admin"
                    name="admin"
                    onChange={(e) => {
                      this.updateRoles(e);
                    }}
                  />
                </div>
              </fieldset>
            </div>
            <div>
              <button
                onClick={() => {
                  this.saveUser();
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

export default withTranslation()(Users);
