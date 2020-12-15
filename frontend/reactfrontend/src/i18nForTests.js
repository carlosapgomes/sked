import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import en from "../public/locales/en/translation.json";
import pt from "../public/locales/pt/translation.json";
i18n.use(initReactI18next).init({
  lng: "pt-br",
  fallbackLng: "en",
  ins: ["translations"],
  defaultNS: "translations",
  debug: true,
  interpolation: {
    escapeValue: false,
  },
  resources: {
    en: { translations: en },
    "pt-BR": { translations: pt },
  },
});
export default i18n;
