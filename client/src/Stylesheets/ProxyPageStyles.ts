import { makeStyles } from '@mui/styles';
import { Theme } from '@mui/material/styles'; // Import Theme

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
  },buttonContainer: {
    margin: '20px 0'
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
  
}));

export default useProxyHostsStyles;
