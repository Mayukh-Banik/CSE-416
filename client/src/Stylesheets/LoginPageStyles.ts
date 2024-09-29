import { makeStyles } from '@mui/styles';
import { Theme } from '@mui/material/styles'; // Import Theme

const useLoginPageStyles = makeStyles((theme: Theme) => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '100vh',
    backgroundColor: theme.palette.background.default,
    padding: theme.spacing(4),
  },
  title: {
    color: theme.palette.primary.main,
    marginBottom: theme.spacing(4),
  },
  loginForm: {
    width: '100%',
    maxWidth: '400px',
    textAlign: 'center',
  },
  textField: {
    marginBottom: theme.spacing(2),
  },
  loginButton: {
    marginTop: theme.spacing(2),
    width: '100%',
  },
  publicKey: {
    marginTop: theme.spacing(4),
    fontWeight: 'bold',
    wordBreak: 'break-word', // Handle long keys
  },
  successMessage: {
    marginTop: theme.spacing(2),
    color: theme.palette.success.main,
  },
}));

export default useLoginPageStyles;
