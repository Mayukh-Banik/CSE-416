import React, { useState } from "react";
import { HDNodeWallet, Mnemonic } from "ethers"; // Import HDNodeWallet and Mnemonic for ethers v6
import { Box, Typography, Button } from "@mui/material";
import * as bip39 from "bip39"; // For generating mnemonic phrases
import { randomBytes } from '@ethersproject/random'; // Correct import for randomBytes
import { hexlify } from '@ethersproject/bytes'; // Correct import for hexlify
import { Buffer } from 'buffer'; // Import buffer polyfill
import useSignUpPageStyles from "../Stylesheets/SignUpPageStyles";

const SignupPage: React.FC = () => {
  const classes = useSignUpPageStyles(); // Custom styles
  const [privateKey, setPrivateKey] = useState<string | null>(null);
  const [publicKey, setPublicKey] = useState<string | null>(null);
  const [mnemonic, setMnemonic] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const generateKeys = () => {
    setErrorMessage(null); // Reset error
    try {
      // Generate random entropy for the mnemonic
      const entropy = randomBytes(16); // 128 bits for 12-word mnemonic
      console.log('Entropy:', entropy); // Log the entropy for debugging

      // Convert Uint8Array to Buffer
      const entropyBuffer = Buffer.from(entropy);
      const mnemonicPhrase = bip39.entropyToMnemonic(entropyBuffer.toString('hex')); // Convert to hex string for bip39
      console.log('Mnemonic Phrase:', mnemonicPhrase); // Log the mnemonic for debugging

      // Convert the mnemonic phrase into the Mnemonic type
      const mnemonicObject = Mnemonic.fromPhrase(mnemonicPhrase); // ethers v6 conversion
      console.log('Mnemonic Object:', mnemonicObject); // Log mnemonic object

      // Create wallet from mnemonic using HDNodeWallet (ethers.js v6)
      const wallet = HDNodeWallet.fromMnemonic(mnemonicObject); 
      console.log('Wallet:', wallet); // Log the generated wallet

      setMnemonic(mnemonicPhrase); // Set the mnemonic phrase
      setPrivateKey(wallet.privateKey); // Private key
      setPublicKey(wallet.address); // Public key
    } catch (error) {
      console.error('Error generating key pair:', error); // Log the error for debugging
      setErrorMessage('Error generating key pair.');
    }
  };

  const handleDownload = () => {
    const element = document.createElement("a");
    const file = new Blob([privateKey as string], { type: "text/plain" });
    element.href = URL.createObjectURL(file);
    element.download = "privateKey.txt";
    document.body.appendChild(element);
    element.click();
  };

  return (
    <Box className={classes.root}>
      <Typography variant="h2" className={classes.title} gutterBottom>
        Set Up Your Account
      </Typography>
      {!privateKey ? (
        <Button variant="contained" color="primary" onClick={generateKeys}>
          Generate Key Pair
        </Button>
      ) : (
        <Box mt={4} textAlign="center">
          <Typography variant="h6" className={classes.publicKey} gutterBottom>
            <strong>Your Public Key:</strong> {publicKey}
          </Typography>
          <Typography variant="body1" className={classes.mnemonic} gutterBottom>
            <strong>12-word Mnemonic Phrase:</strong> {mnemonic}
          </Typography>
          <Button
            variant="outlined"
            color="primary"
            onClick={handleDownload}
            className={classes.button}
          >
            Download Private Key
          </Button>
          <Typography variant="body1" className={classes.warning}>
            Warning: If you lose your private key or mnemonic, you will lose
            access to your account.
          </Typography>
        </Box>
      )}
      {errorMessage && (
        <Typography color="error" variant="body2" mt={2}>
          {errorMessage}
        </Typography>
      )}
    </Box>
  );
};

export default SignupPage;
