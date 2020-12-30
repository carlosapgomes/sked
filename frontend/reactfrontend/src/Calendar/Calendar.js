import React, { Component } from "react";
import { withTranslation } from "react-i18next";
import "./Calendar.css";

class Calendar extends Component {
  constructor(props) {
    super(props);

    this.state = {};
  }

  render() {
    const { t } = this.props;
    return (
      <div className="Calendar">
        <div className="Labels">
          <span>{t("Sunday")}</span>
          <span>{t("Monday")}</span>
          <span>{t("Tuesday")}</span>
          <span>{t("Wednesday")}</span>
          <span>{t("Thursday")}</span>
          <span>{t("Friday")}</span>
          <span>{t("Saturday")}</span>
        </div>
        <div className="Dates">
          {this.props.view.map((w) => {
            let i = 0;
            return w.map((d) => {
              return <span key={i++}>{d > 0 ? d : ""}</span>;
            });
          })}
        </div>
      </div>
    );
  }
}
export default withTranslation()(Calendar);
