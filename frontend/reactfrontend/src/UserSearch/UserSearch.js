import React, { Component } from "react";
import { withTranslation } from "react-i18next";

class UserSearch extends Component {
  constructor(props) {
    super(props);

    this.state = {
      searchField: "",
      userSearchResult: [],
      selectedUserValue: "selectAnOption",
      selectedUser: null,
    };
  }
  resetInitialState() {
    this.setState({
      searchField: "",
      userSearchResult: [],
      selectedUserValue: "selectAnOption",
      selectedUser: null,
    });
  }
  setSelectedUser(e) {
    let idx = e.target.selectedIndex - 1;
    this.setState({
      selectedUserValue: this.state.userSearchResult[idx].id,
      selectedUser: { ...this.state.userSearchResult[idx] },
      searchField: "",
    });
    // update selectedUser on parent component
    this.props.setSelectedUser({ ...this.state.userSearchResult[idx] });
  }
  setSearchField(s) {
    this.setState({
      searchField: s,
    });
  }
  searchUser() {
    if (!this.state.searchField || this.state.searchField.length < 3) {
      return;
    }
    let str = this.state.searchField.trim().split(/\s+/).join(" ");
    if (str.length < 3) {
      return;
    }
    let ajax = new XMLHttpRequest();
    let url = "https://dev.local/api/users?name=" + str;
    ajax.open("GET", url, true);
    ajax.withCredentials = true;
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send();
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        if (!ajax.responseText) {
          window.alert("Could not find any user");
        } else {
          let data = JSON.parse(ajax.responseText);
          if (data) {
            this.setState({
              userSearchResult: [...data],
              selectedUserValue: "selectAnOption",
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
        <label htmlFor="pctsearch">{t("user")}: </label>
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
            this.searchUser();
          }}
        >
          {t("Search")}
        </button>
        <div hidden={this.state.userSearchResult.length <= 0}>
          <select
            name="searchresult"
            id="searchresult"
            onChange={(e) => {
              this.setSelectedUser(e);
            }}
            value={this.state.selectedUserValue}
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
            {this.state.userSearchResult.map((p) => {
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
export default withTranslation()(UserSearch);
