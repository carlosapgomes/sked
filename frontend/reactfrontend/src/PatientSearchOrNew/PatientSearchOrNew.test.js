import React from "react";
import { render } from "@testing-library/react";
import PatientSearchOrNew from "./PatientSearchOrNew";

it("renders without crashing", () => {
  const div = document.createElement("div");
  render(<PatientSearchOrNew currentPatient={null} />, div);
});
