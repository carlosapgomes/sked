import React, { Component } from "react";

export default class Calendar extends Component {
  constructor(props) {
    super(props);

    this.state = {};
  }

  render() {
    return (
      <div className="Calendar">
        <div className="Labels"></div>
        <div className="dates"></div>
      </div>
    );
  }
}
