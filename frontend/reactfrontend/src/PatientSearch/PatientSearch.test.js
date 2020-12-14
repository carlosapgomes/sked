import React from "react";
import { render } from "@testing-library/react";
import PatientSearch from "./PatientSearch";

it("renders without crashing", () => {
  const div = document.createElement("div");
  render(<PatientSearch currentPatient={null} />, div);
});
