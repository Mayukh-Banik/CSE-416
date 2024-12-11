import React, { useState, useEffect, useCallback } from 'react';
import {
  Box, Table, TableBody, TableCell, TableHead, TableRow, TableContainer,
  Typography, Button, Paper,
  TextField, Grid, Alert, CircularProgress, IconButton, InputAdornment,
  Snackbar,
} from '@mui/material';
import Sidebar from './Sidebar';
import { useTheme } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Transaction } from '../models/transactions';
import { UnspentTransaction } from '../models/unspentTransaction';
import { Visibility, VisibilityOff, ContentCopy, Refresh } from '@mui/icons-material';
import { CopyToClipboard } from 'react-copy-to-clipboard';
import MuiAlert from '@mui/material/Alert';

// Snackbar Alert Component
const AlertComponent = React.forwardRef<HTMLDivElement, any>(function Alert(
  props,
  ref,
) {
  return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

// Set up API base URL (environment variables recommended)
const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080';

const GlobalTransactions: React.FC = () => {
  const theme = useTheme();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [unspentTxs, setUnspentTxs] = useState<UnspentTransaction[]>([]);
  const [currentAddress, setCurrentAddress] = useState<string>('');
  const navigate = useNavigate();

  const [passphrase, setPassphrase] = useState<string>('');
  const [txid, setTxid] = useState<string>('');
  const [dst, setDst] = useState<string>('');
  const [amount, setAmount] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');
  const [showPassphrase, setShowPassphrase] = useState<boolean>(false);

  // Snackbar state
  const [openSnackbar, setOpenSnackbar] = useState<boolean>(false);
  const [snackbarMessage, setSnackbarMessage] = useState<string>('');
  const [snackbarSeverity, setSnackbarSeverity] = useState<'success' | 'error'>('success');

  // Snackbar close handler
  const handleCloseSnackbar = (
    event?: React.SyntheticEvent | Event,
    reason?: string
  ) => {
    if (reason === 'clickaway') {
      return;
    }
    setOpenSnackbar(false);
  };

  // Fetch unspent transactions data
  const fetchUnspentTransactions = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/btc/listunspent`);
      const data = await response.json();
      if (data.status === 'success' && Array.isArray(data.data)) {
        console.log("Unspent Transactions: ", data.data);
        setUnspentTxs(data.data);
      } else {
        console.error("Unexpected response format for unspent transactions:", data);
        setError("Failed to load unspent transaction data.");
      }
    } catch (error) {
      console.error("Failed to fetch unspent transactions: ", error);
      setError("Failed to load unspent transaction data.");
    }
  }, []);

  // Fetch current address
  const fetchCurrentAddress = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/btc/currentaddress`);
      const data = await response.json();
      if (data.status === 'success' && typeof data.data === 'string') {
        setCurrentAddress(data.data);
      } else {
        console.error("Unexpected response format for current address:", data);
        setError("Failed to load the current address.");
      }
    } catch (error) {
      console.error("Failed to fetch current address:", error);
      setError("Failed to load the current address.");
    }
  }, []);

  // Fetch initial data
  useEffect(() => {
    fetchUnspentTransactions();
    fetchCurrentAddress();
  }, [fetchUnspentTransactions, fetchCurrentAddress]);

  const handleTransactionSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');

    // Example input validation
    if (amount <= 0) {
      setError('Amount must be greater than 0.');
      setLoading(false);
      return;
    }

    const txidPattern = /^[A-Fa-f0-9]{64}$/;
    if (!txidPattern.test(txid)) {
      setError('Invalid Transaction ID format.');
      setLoading(false);
      return;
    }

    const addressPattern = /^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$/;
    if (!addressPattern.test(dst)) {
      setError('Invalid Destination Address format.');
      setLoading(false);
      return;
    }

    try {
      const response = await axios.post(`${API_BASE_URL}/api/btc/transaction`, {
        passphrase,
        txid,
        dst,
        amount,
      });

      if (response.data.status === 'success') {
        setSuccess(`Transaction completed successfully. TxID: ${response.data.data.txid}`);

        // Refresh transaction list
        // await fetchTransactions();
        // Refresh unspent transaction list
        await fetchUnspentTransactions();
        setPassphrase('');
        setTxid('');
        setDst('');
        setAmount(0);
      } else {
        setError(response.data.message || 'Transaction failed.');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'An error occurred during the transaction.');
    } finally {
      setLoading(false);
    }
  };

  const handleRefresh = () => {
    setLoading(true);
    setError('');
    setSuccess('');
    fetchUnspentTransactions().then(() => {
      fetchCurrentAddress().finally(() => setLoading(false));
    });
  };

  const renderUnspentTransactionsTable = () => {
    // Filter unspent transactions matching the current address
    const filteredUnspentTxs = unspentTxs.filter(tx => tx.address === currentAddress);

    if (filteredUnspentTxs.length === 0) {
      return (
        <Paper sx={{ marginBottom: 4, padding: 2 }}>
          {/* Display current address */}
          <Typography variant="h6" gutterBottom>
            Current Address: {currentAddress}
          </Typography>

          {/* Empty State Message */}
          <Box sx={{ textAlign: 'center', padding: 4 }}>
            <Typography variant="body1" gutterBottom>
              No unspent transactions available.
            </Typography>
            {/* Optional: Add an illustration or icon for empty state */}
            <img
              src="/mining.gif" // 수정된 이미지 경로
              alt="No Data"
              style={{ marginTop: 16, width: '150px', height: '150px' }}
            />
          </Box>
        </Paper>
      );
    }

    return (
      <Paper sx={{ marginBottom: 4, padding: 2 }}>
        {/* Display current address and Refresh button */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 2 }}>
          <Typography variant="h6">
            Current Address: {currentAddress}
          </Typography>
          <Button
            variant="outlined"
            startIcon={<Refresh />}
            onClick={handleRefresh}
            disabled={loading}
          >
            {loading ? <CircularProgress size={20} /> : 'Refresh'}
          </Button>
        </Box>

        {/* Copy all data button */}
        <Box sx={{ padding: 2, display: 'flex', justifyContent: 'flex-end' }}>
          <CopyToClipboard
            text={filteredUnspentTxs.map(tx => `TxID: ${tx.txid}, Amount: ${tx.amount}, Confirmations: ${tx.confirmations}`).join('\n')}
            onCopy={() => {
              setSnackbarMessage('All data copied successfully!');
              setSnackbarSeverity('success');
              setOpenSnackbar(true);
            }}
          >
            <Button variant="contained" color="secondary" startIcon={<ContentCopy />}>
              Copy All Data
            </Button>
          </CopyToClipboard>
        </Box>

        {/* Responsive Table Container */}
        <TableContainer sx={{ maxHeight: 440, overflowX: 'auto' }}>
          <Table stickyHeader aria-label="unspent transactions table">
            <TableHead>
              <TableRow>
                <TableCell>Transaction ID</TableCell>
                <TableCell>Amount</TableCell>
                <TableCell>Confirmations</TableCell>
                <TableCell>Copy</TableCell> {/* Copy button column */}
              </TableRow>
            </TableHead>
            <TableBody>
              {filteredUnspentTxs.map((tx, index) => (
                <TableRow key={index}>
                  <TableCell>{tx.txid}</TableCell>
                  <TableCell>{tx.amount}</TableCell>
                  <TableCell>{tx.confirmations}</TableCell>
                  <TableCell>
                    <CopyToClipboard
                      text={tx.txid} // Copy only txid
                      onCopy={() => {
                        setSnackbarMessage('TxID copied successfully!');
                        setSnackbarSeverity('success');
                        setOpenSnackbar(true);
                      }}
                    >
                      <IconButton aria-label="copy">
                        <ContentCopy />
                      </IconButton>
                    </CopyToClipboard>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
    );
  };

  return (
    <Box
      sx={{
        padding: 2,
        marginTop: '70px',
        marginLeft: `${drawerWidth}px`,
        [theme.breakpoints.down('sm')]: {
          marginLeft: `${collapsedDrawerWidth}px`,
        },
      }}
    >
      <Sidebar />
      <Typography variant="h4" component="h1" gutterBottom>
        Transaction
      </Typography>

      {/* Transaction Form */}
      <Paper sx={{ padding: 3, marginBottom: 4 }}>
        <Typography variant="h6" gutterBottom>
          Perform a Transaction
        </Typography>
        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}
        <form onSubmit={handleTransactionSubmit}>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <TextField
                label="Passphrase"
                variant="outlined"
                fullWidth
                required
                type={showPassphrase ? "text" : "password"}
                value={passphrase}
                onChange={(e) => setPassphrase(e.target.value)}
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        aria-label="toggle passphrase visibility"
                        onClick={() => setShowPassphrase(!showPassphrase)}
                        edge="end"
                      >
                        {showPassphrase ? <VisibilityOff /> : <Visibility />}
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                label="Transaction ID (TxID)"
                variant="outlined"
                fullWidth
                required
                value={txid}
                onChange={(e) => setTxid(e.target.value)}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                label="Destination Address"
                variant="outlined"
                fullWidth
                required
                value={dst}
                onChange={(e) => setDst(e.target.value)}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                label="Amount"
                variant="outlined"
                type="number"
                inputProps={{ step: "0.00000001" }}
                fullWidth
                required
                value={amount}
                onChange={(e) => setAmount(parseFloat(e.target.value))}
              />
            </Grid>
            <Grid item xs={12}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                disabled={loading}
                startIcon={loading && <CircularProgress size={20} />}
              >
                {loading ? 'Processing...' : 'Submit Transaction'}
              </Button>
            </Grid>
          </Grid>
        </form>
      </Paper>

      {/* Render Unspent Transactions Table or Empty State */}
      {unspentTxs.length > 0 ? (
        renderUnspentTransactionsTable()
      ) : (
        // Even if unspentTxs is empty, currentAddress should still be displayed
        <Paper sx={{ marginBottom: 4, padding: 2 }}>
          <Typography variant="h6" gutterBottom>
            Current Address: {currentAddress}
          </Typography>
          <Box sx={{ textAlign: 'center', padding: 4 }}>
            <Typography variant="body1" gutterBottom>
              No unspent transactions available.
            </Typography>
            {/* Optional: Add an illustration or icon for empty state */}
            {/* Refresh */}
            <Box sx={{ marginTop: 2 }}>
              <Button
                variant="outlined"
                startIcon={<Refresh />}
                onClick={handleRefresh}
                disabled={loading}
              >
                {loading ? <CircularProgress size={20} /> : 'Refresh'}
              </Button>
            </Box>
          </Box>
        </Paper>
      )}

      {/* List of existing transactions */}
      {/* {renderTransactionsTable(transactions)} */}

      {/* Snackbar Notifications */}
      <Snackbar
        open={openSnackbar}
        autoHideDuration={3000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <AlertComponent
          onClose={handleCloseSnackbar}
          severity={snackbarSeverity}
          sx={{ width: '100%' }}
        >
          {snackbarMessage}
        </AlertComponent>
      </Snackbar>
    </Box>
  );
};

export default GlobalTransactions;
