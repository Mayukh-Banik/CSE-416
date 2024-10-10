import React, { useState } from 'react';
import * as bitcoin from 'bitcoinjs-lib';
import * as ecc from 'tiny-secp256k1';
import { ECPairFactory } from 'ecpair';

const ECPair = ECPairFactory(ecc);

const SignUpPage: React.FC = () => {
  const [walletAddress, setWalletAddress] = useState<string | null>(null);
  const [privateKey, setPrivateKey] = useState<string | null>(null);
  const [isSubmitted, setIsSubmitted] = useState(false);

  const generateWallet = async () => {
    // Generate a Bitcoin key pair using ECPair
    const keyPair = ECPair.makeRandom();
    const { address } = bitcoin.payments.p2pkh({ pubkey: keyPair.publicKey });
    const privateKey = keyPair.toWIF();

    // Set the state to show the address and private key
    setWalletAddress(address || null);
    setPrivateKey(privateKey);
    setIsSubmitted(true);

    // Send the public key (walletAddress) to the backend to register
    try {
      await fetch('/api/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ walletAddress: address }),
      });
    } catch (error) {
      console.error('Error registering wallet:', error);
    }
  };

  return (
    <div>
      <h1>Sign Up</h1>
      {!isSubmitted ? (
        <button onClick={generateWallet}>Generate Wallet</button>
      ) : (
        <div>
          <p>Your wallet address (public key): {walletAddress}</p>
          <p>Your private key: {privateKey}</p>
          <p>
            <strong>Important:</strong> Keep your private key secure and do not share it with
            anyone.
          </p>
        </div>
      )}
    </div>
  );
}

export default SignUpPage;