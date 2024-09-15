import { makeStyles } from "@mui/material";

const useStyles = makeStyles({
  container: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
    height: '80vh', // Vertically center the content
    textAlign: 'center',
  },
  heading: {
    fontSize: '3rem',
    fontWeight: 700,
    color: '#444',
    marginBottom: '2rem',
  },
  subtext: {
    fontSize: '1.2rem',
    color: '#777',
    marginBottom: '4rem',
  },
  button: {
    width: '100%',
    padding: '15px 0',
    fontSize: '1.2rem',
    borderRadius: '8px',
    boxShadow: '0px 4px 12px rgba(0, 0, 0, 0.1)',
    transition: 'all 0.3s ease',
    marginBottom: '3rem', // Spacing between buttons
    '&:hover': {
      backgroundColor: '#95bfe7',
    },
  },
});

export default useStyles;
