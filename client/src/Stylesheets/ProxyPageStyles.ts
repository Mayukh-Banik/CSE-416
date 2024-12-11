import { makeStyles } from '@mui/styles';
import { Theme } from '@mui/material/styles';
const drawerWidth = 300;
const collapsedDrawerWidth = 100;
const useProxyHostsStyles = makeStyles((theme: Theme) => ({
  container: {
    padding: '20px',
    fontFamily: 'Arial, sans-serif',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '20px',
  },
  addButton: {
    padding: '10px 15px',
    backgroundColor: '#007bff',
    color: 'white',
    border: 'none',
    borderRadius: '5px',
    cursor: 'pointer',
  },
  table: {
    width: '100%',
    borderCollapse: 'collapse',
    boxShadow: '0 0 10px rgba(0, 0, 0, 0.1)',
  },
  sourceColumn: {
    display: 'flex',
    alignItems: 'left',
  },
  avatar: {
    borderRadius: '50%',
    marginRight: '10px',
  },
  createdDate: {
    fontSize: '12px',
    color: '#888',
  },
  statusColumn: {
    display: 'flex',
    alignItems: 'center',
  },
  statusIndicator: {
    width: '10px',
    height: '10px',
    borderRadius: '50%',
    marginRight: '5px',
  },
  buttonContainer: {
    margin: '20px 0',
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
    gap: '10px',
    marginBottom: '20px',
  },
  submitButton: {
    backgroundColor: 'green',
    color: 'white',
    padding: '10px',
    cursor: 'pointer',
    border: 'none',
  },
  proxyTable: {
    marginTop: 2,
  },
  proxyButton: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: 1,
  },
  popup: {
    position: 'fixed',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    backgroundColor: theme.palette.background.paper,
    padding: theme.spacing(2),
    boxShadow: theme.shadows[5],
    borderRadius: theme.shape.borderRadius,
    zIndex: 1000,
  },
  historyContainer: {
    marginTop: 2,
  },
  historyTable: {
    marginTop: 2,
  },
  boxContainer: {
    padding: '16px', 
    marginTop: '70px',
    marginLeft: `${drawerWidth}px`,
    transition: 'margin-left 0.3s ease',
    [theme.breakpoints.down('sm')]: {
      marginLeft: `${collapsedDrawerWidth}px`,
    },
  },
}));

export default useProxyHostsStyles;