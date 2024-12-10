// src/models/unspentTransaction.ts

export interface UnspentTransaction {
    txid: string;
    vout: number;
    address: string;
    scriptPubKey: string;
    amount: number;
    confirmations: number;
  }
  