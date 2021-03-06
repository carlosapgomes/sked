import React, { Component } from "react";
import { withTranslation } from "react-i18next";
import UserSearch from "../UserSearch/UserSearch";
class Users extends Component {
  constructor(props) {
    super(props);

    this.state = {
      id: "",
      name: "",
      email: "",
      phone: "",
      roles: [],
      showUpdateButton: false,
      isAdmin: false,
      isClerk: false,
      isDoctor: false,
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
      this.state.phone === "" ||
      this.state.roles.length === 0
    ) {
      window.alert("Please, fill the requested data");
      return;
    }
    let newUser = {
      Name: this.state.name,
      Email: this.state.email,
      Phone: this.state.phone,
      Roles: [...this.state.roles],
    };
    console.log(newUser);
    let ajax = new XMLHttpRequest();
    let url = "/api/users";
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
  updateUser() {
    if (
      this.state.id === "" ||
      this.state.name === "" ||
      this.state.email === "" ||
      this.state.phone === "" ||
      this.state.roles.length === 0
    ) {
      window.alert("Please, fill the requested data");
      return;
    }
    let updatedUser = {
      ID: this.state.id,
      Name: this.state.name,
      Email: this.state.email,
      Phone: this.state.phone,
      Roles: [...this.state.roles],
    };
    console.log(updatedUser);
    let ajax = new XMLHttpRequest();
    let url = "/api/users";
    ajax.open("PUT", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify(updatedUser));
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        window.alert("User updated");
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
      email: "",
      phone: "",
      roles: [],
      isAdmin: false,
      isClerk: false,
      isDoctor: false,
    });
  }
  localSubmitHandler(e) {
    e.preventDefault();
  }
  updateRoles(r) {
    let roles = [...this.state.roles];
    let i = roles.indexOf(r);
    if (i < 0) {
      roles.push(r);
    } else {
      roles.splice(i, 1);
    }
    this.setState({
      roles: [...roles],
    });
    this.updateRolesCheckBoxes(roles);
  }
  updateRolesCheckBoxes(roles) {
    let isClerk,
      isDoctor,
      isAdmin = false;
    if (roles.indexOf("Clerk") > -1) {
      isClerk = true;
    }

    if (roles.indexOf("Doctor") > -1) {
      isDoctor = true;
    }
    if (roles.indexOf("Admin") > -1) {
      isAdmin = true;
    }
    this.setState({
      isClerk: isClerk,
      isDoctor: isDoctor,
      isAdmin: isAdmin,
    });
  }

  setSelectedUser(u) {
    if (!u) {
      this.clearForm();
    } else {
      this.setState({
        id: u.id,
        name: u.name,
        email: u.email,
        phone: u.phone,
        roles: [...u.roles],
        showUpdateButton: true,
      });
    }
    this.updateRolesCheckBoxes(u.roles);
  }
  render() {
    const { t } = this.props;
    return (
      <div>
        <h1>{t("Users")}</h1>
        <section>
          <UserSearch
            setSelectedUser={(u) => {
              this.setSelectedUser(u);
            }}
          />
          <br />
        </section>
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
                disabled={this.state.showUpdateButton}
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
                disabled={this.state.showUpdateButton}
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
              <fieldset>
                <legend>{t("Roles")}:</legend>
                <div>
                  <label htmlFor="clerk">{t("Clerk")}</label>
                  <input
                    type="checkbox"
                    name="clerk"
                    checked={this.state.isClerk || false}
                    onChange={() => {
                      this.updateRoles("Clerk");
                    }}
                  />
                </div>
                <div>
                  <label htmlFor="clerk">{t("Doctor")}</label>
                  <input
                    type="checkbox"
                    name="doctor"
                    checked={this.state.isDoctor || false}
                    onChange={() => {
                      this.updateRoles("Doctor");
                    }}
                  />
                </div>
                <div>
                  <label htmlFor="clerk">{t("Admin")}</label>
                  <input
                    type="checkbox"
                    name="admin"
                    checked={this.state.isAdmin || false}
                    onChange={() => {
                      this.updateRoles("Admin");
                    }}
                  />
                </div>
              </fieldset>
            </div>
            <div>
              <button
                hidden={this.state.showUpdateButton}
                onClick={() => {
                  this.saveUser();
                }}
              >
                {t("Save")}
              </button>
              <button
                hidden={!this.state.showUpdateButton}
                onClick={() => {
                  this.updateUser();
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
          </form>
        </section>
      </div>
    );
  }
}

export default withTranslation()(Users);
