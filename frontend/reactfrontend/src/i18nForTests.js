import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import en from "../public/locales/en/translation.json";

i18n.use(initReactI18next).init({
  lng: "en",
  fallbackLng: "en",
  ins: ["translations"],
  defaultNS: "translations",

  debug: true,
  interpolation: {
    escapeValue: false,
  },
  resources: {
    en,
  },
});
export default i18n;
