import React, { useState } from "react";
import { HDNodeWallet, Mnemonic } from "ethers"; // Import HDNodeWallet and Mnemonic from ethers.js v6
import { Box, Typography, Button, TextField } from "@mui/material";
import useLoginPageStyles from "../Stylesheets/LoginPageStyles"; // Custom styles

const LoginPage2: React.FC = () => {
  const classes = useLoginPageStyles(); // Custom styles
  const [mnemonicInput, setMnemonicInput] = useState<string>(""); // Input for mnemonic
  const [publicKey, setPublicKey] = useState<string | null>(null); // Store public key
  const [errorMessage, setErrorMessage] = useState<string | null>(null); // Error handling

  const handleLogin = () => {
    setErrorMessage(null); // Reset error message
    try {
      // Validate the input is a valid mnemonic
      const mnemonicObject = Mnemonic.fromPhrase(mnemonicInput); // Convert to Mnemonic type
      
      // Recover wallet from mnemonic
      const wallet = HDNodeWallet.fromMnemonic(mnemonicObject);
      
      setPublicKey(wallet.address); // Set the public key
    } catch (error) {
      console.error('Login error:', error);
      setErrorMessage('Invalid mnemonic phrase. Please check and try again.');
    }
  };

  return (
    <Box className={classes.root}>
      <Typography variant="h2" className={classes.title} gutterBottom>
        Log In
      </Typography>
      {!publicKey ? (
        <Box className={classes.loginForm}>
          <TextField
            label="Enter 12-word Mnemonic Phrase"
            variant="outlined"
            fullWidth
            multiline
            rows={3}
            value={mnemonicInput}
            onChange={(e) => setMnemonicInput(e.target.value)}
            className={classes.textField}
          />
          <Button
            variant="contained"
            color="primary"
            onClick={handleLogin}
            className={classes.loginButton}
          >
            Log In
          </Button>
          {errorMessage && (
            <Typography color="error" variant="body2" mt={2}>
              {errorMessage}
            </Typography>
          )}
        </Box>
      ) : (
        <Box mt={4} textAlign="center">
          <Typography variant="h6" className={classes.publicKey}>
            <strong>Your Public Key:</strong> {publicKey}
          </Typography>
          <Typography variant="body1" className={classes.successMessage}>
            You have successfully logged in using your mnemonic phrase.
          </Typography>
        </Box>
      )}
    </Box>
  );
};

export default LoginPage2;
