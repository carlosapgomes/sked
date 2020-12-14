import i18n from "i18next";
import { initReactI18next } from "react-i18next";

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
    en: {},
    pt: {},
  },
});
export default i18n;
