import { makeStyles } from "@mui/styles";


const useRegisterPageStyles = makeStyles({
  container: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
    height: '80vh',
    textAlign: 'center',
  },
  icon: {
    fontSize: '3rem',
    color: '#4285f4', // Google-like blue for the icon
    marginBottom: '1rem',
  },
  form: {
    width: '100%', // Make form take full width
    maxWidth: '400px', // Optional: Add max width
  },
  inputField: {
    width: '100%',
    marginBottom: '1.5rem',
  },
  button: {
    width: '100%',
    padding: '12px 0',
    fontSize: '1.2rem',
    backgroundColor: 'primary', // Primary color from theme
    color: '#fff',
    borderRadius: '8px',
    boxShadow: '0px 4px 12px rgba(0, 0, 0, 0.1)',
    marginBottom: '1.5rem',
    '&:hover': {
      backgroundColor: '#357ae8', // Darker on hover
    },
  },
  link: {
    color: '#4285f4',
    textDecoration: 'none',
    fontSize: '0.9rem',
    '&:hover': {
      textDecoration: 'underline',
    },
  },
});

export default useRegisterPageStyles;
