import { makeStyles } from "@mui/styles";
import { Theme } from "@mui/material/styles"; // Import Theme

// Define custom styles for SignUpPage
const useSignUpPageStyles = makeStyles((theme: Theme) => ({
  root: {
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    alignItems: "center",
    minHeight: "100vh",
    backgroundColor: theme.palette.background.default,
    padding: theme.spacing(4),
  },
  title: {
    color: theme.palette.primary.main,
    marginBottom: theme.spacing(2),
  },
  publicKey: {
    marginTop: theme.spacing(4),
    wordBreak: "break-word", // Handle long public keys
    fontWeight: "bold",
  },
  mnemonic: {
    fontSize: "1.1rem",
    fontStyle: "italic",
    color: theme.palette.secondary.main,
    marginTop: theme.spacing(2),
  },
  button: {
    marginTop: theme.spacing(2),
  },
  warning: {
    marginTop: theme.spacing(2),
    color: theme.palette.error.main,
  },
}));

export default useSignUpPageStyles;
