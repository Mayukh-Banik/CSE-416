declare module 'tiny-secp256k1';

// global.d.ts
declare module 'crypto-browserify' {
    import { Hash, Hmac, Sign, Verify } from 'crypto';
  
    export function createHash(algorithm: string): Hash;
    export function createHmac(algorithm: string, key: string | Buffer): Hmac;
    export function createSign(algorithm: string): Sign;
    export function createVerify(algorithm: string): Verify;
    export function randomBytes(size: number): Buffer;
  }
  