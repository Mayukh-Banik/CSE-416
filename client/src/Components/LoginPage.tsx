import React, { useState } from 'react';
import { TextField, Button, Container, Typography, Link, Box } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import useRegisterPageStyles from '../Stylesheets/RegisterPageStyles';
import Header from './Header';
import { useNavigate } from 'react-router-dom';

const PORT = 8080;

const LoginPage: React.FC = () => {
    const classes = useRegisterPageStyles();
    const navigate = useNavigate();

    const [publicKey, setPublicKey] = useState('');
    const [privateKey, setPrivateKey] = useState(''); // 프라이빗 키 상태 추가
    const [error, setError] = useState<string | null>(null);
    const [challenge, setChallenge] = useState<string | null>(null);

    // 폼 유효성 검사 (퍼블릭 키만 확인)
    const validateForm = () => {
        if (!publicKey) {
            setError('Public key is required.');
            return false;
        }
        setError(null);
        return true;
    };

    // 서버에 퍼블릭 키를 보내고 챌린지를 받는 함수
    const handleLogin = async () => {
        if (!validateForm()) return;

        const requestBody = JSON.stringify({ public_key: publicKey });
        console.log("Request Body:", requestBody);

        try {
            const response = await fetch(`http://localhost:${PORT}/api/auth/request-challenge`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: requestBody
            });

            if (!response.ok) {
                const errorMessage = await response.text();
                console.error("Server Error:", errorMessage);
                throw new Error('Failed to request challenge');
            }

            const data = await response.json();
            const receivedChallenge = data.challenge;
            setChallenge(receivedChallenge);

            console.log("Challenge received:", receivedChallenge);
        } catch (err) {
            console.error("Error requesting challenge:", err);
            setError('Failed to request challenge.');
        }
    };

    // 서명 처리 함수
    const handleSign = async () => {
        if (!challenge || !privateKey) {
            setError('Challenge and private key are required.');
            console.log('Error: Missing challenge or private key');
            return;
        }

        try {
            console.log('Converting PEM private key to ArrayBuffer...');
            const privateKeyBuffer = pemToArrayBuffer(privateKey); // 올바르게 변환된 ArrayBuffer
            console.log('Private key ArrayBuffer:', privateKeyBuffer);

            // Web Crypto API를 사용하여 프라이빗 키를 가져옴
            console.log('Importing private key...');
            const privateKeyObj = await window.crypto.subtle.importKey(
                "pkcs8", // PKCS#8 형식
                privateKeyBuffer,
                {
                    name: "RSASSA-PKCS1-v1_5",
                    hash: { name: "SHA-256" },
                },
                true,
                ["sign"]
            );
            console.log('Private key imported successfully:', privateKeyObj);

            // 챌린지 데이터를 바이트 배열로 변환
            console.log('Encoding challenge to byte array...');
            const encoder = new TextEncoder();
            const challengeData = encoder.encode(challenge);
            console.log('Challenge byte array:', challengeData);

            // Web Crypto API를 사용하여 서명 생성
            console.log('Signing challenge...');
            const signatureBuffer = await window.crypto.subtle.sign(
                "RSASSA-PKCS1-v1_5", // RSA 서명 알고리즘
                privateKeyObj,        // 가져온 프라이빗 키
                challengeData         // 서명할 데이터 (챌린지)
            );

            // 서명을 Hex 형식으로 변환
            const signature = Array.from(new Uint8Array(signatureBuffer))
                .map(b => b.toString(16).padStart(2, "0"))
                .join("");
            console.log('Signature generated:', signature);

            // 서명을 서버로 전송할 준비
            const requestBody = JSON.stringify({ public_key: publicKey, signature });
            console.log('Signing Request Body:', requestBody);

            const response = await fetch(`http://localhost:${PORT}/api/auth/verify-challenge`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: requestBody
            });

            if (!response.ok) {
                const errorMessage = await response.text();
                console.error('Server Error:', errorMessage);
                throw new Error('Failed to verify challenge');
            }

            const data = await response.json();
            console.log('Verification Response:', data);
        } catch (err) {
            console.error('Error signing challenge:', err);
            setError('Failed to sign challenge.');
        }
    };

    // PEM 형식의 RSA 프라이빗 키를 ArrayBuffer로 변환하는 함수
    const pemToArrayBuffer0 = (pem: string): ArrayBuffer => {
        const pemHeader = "-----BEGIN RSA PRIVATE KEY-----";
        const pemFooter = "-----END RSA PRIVATE KEY-----";

        // PEM 헤더 및 푸터 제거
        let pemContents = pem.replace(pemHeader, "").replace(pemFooter, "");

        // 공백 및 줄바꿈 제거
        pemContents = pemContents.replace(/\s+/g, "");

        try {
            return base64ToArrayBuffer(pemContents);
        } catch (error) {
            console.error("Error converting PEM to ArrayBuffer:", error);
            throw new Error("Invalid PEM format or Base64 encoding");
        }
    };

    const pemToArrayBuffer = (pem: string): ArrayBuffer => {
        const pemHeader = "-----BEGIN PRIVATE KEY-----";  // Adjust this if necessary
        const pemFooter = "-----END PRIVATE KEY-----";
    
        // Check if the headers and footers are present
        if (!pem.includes(pemHeader) || !pem.includes(pemFooter)) {
            console.error("Invalid PEM format: Missing headers or footers.");
            throw new Error("Invalid PEM format.");
        }
    
        // Remove headers, footers, and whitespace
        let pemContents = pem.replace(pemHeader, "").replace(pemFooter, "").replace(/\s+/g, "");
    
        // Log the content to verify
        console.log("Cleaned PEM Content for Base64 decoding:", pemContents);
    
        // Convert to ArrayBuffer
        return base64ToArrayBuffer(pemContents);
    };
    // Base64 문자열을 ArrayBuffer로 변환하는 함수
    const base64ToArrayBuffer = (base64: string): ArrayBuffer => {
        try {
            const binaryString = window.atob(base64);  // Base64 디코딩
            const len = binaryString.length;
            const bytes = new Uint8Array(len);
            for (let i = 0; i < len; i++) {
                bytes[i] = binaryString.charCodeAt(i);
            }
            return bytes.buffer;
        } catch (error) {
            console.error("Error decoding Base64 string:", error);
            throw new Error("Invalid Base64 string");
        }
    };
    
    const handleSignup = () => {
        navigate('/');
    };

    return (
        <>
            <Header />
            <Container
                maxWidth="sm"
                sx={{
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "center",
                    height: "80vh",
                    textAlign: "center",
                    marginTop: "2rem",
                }}
            >
                <LockOutlinedIcon sx={{ color: "primary.main", fontSize: "3rem", marginBottom: "0.5rem" }} />
                <Typography variant="h4" sx={{ fontWeight: "normal", mb: 2 }}>
                    Log In
                </Typography>

                {error && (
                    <Typography color="error" sx={{ marginBottom: '1rem' }}>
                        {error}
                    </Typography>
                )}

                <Box
                    component="form"
                    sx={{
                        display: "flex",
                        flexDirection: "column",
                        gap: 2,
                        width: "100%",
                        maxWidth: "450px",
                        padding: "2rem",
                        borderRadius: "12px",
                        boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
                    }}
                >
                    <TextField
                        label="Public Key"
                        type="text"
                        variant="outlined"
                        fullWidth
                        required
                        value={publicKey}
                        onChange={(e) => {
                            console.log("Public Key Input:", e.target.value);
                            setPublicKey(e.target.value);
                        }}
                    />
                    <Button
                        variant="contained"
                        color="primary"
                        sx={{
                            padding: "12px 0",
                            fontSize: "1.2rem",
                            borderRadius: "8px",
                            boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
                            "&:hover": {
                                backgroundColor: "#1976d2",
                            },
                        }}
                        onClick={handleLogin}
                    >
                        Request Challenge
                    </Button>
                </Box>

                {challenge && ( // 챌린지가 있을 때 프라이빗 키 입력 폼을 보여줌
                    <Box
                        component="form"
                        sx={{
                            display: "flex",
                            flexDirection: "column",
                            gap: 2,
                            width: "100%",
                            maxWidth: "450px",
                            padding: "2rem",
                            borderRadius: "12px",
                            boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
                        }}
                    >
                        <TextField
                            label="Private Key"
                            type="password"
                            variant="outlined"
                            fullWidth
                            required
                            value={privateKey}
                            onChange={(e) => {
                                console.log("Private Key Input:", e.target.value);
                                setPrivateKey(e.target.value);
                            }}
                        />
                        <Button
                            variant="contained"
                            color="primary"
                            onClick={handleSign} // 서명 버튼 클릭 시
                        >
                            Sign Challenge
                        </Button>
                    </Box>
                )}

                <Typography sx={{ marginTop: "1rem" }}>
                    <Link
                        onClick={handleSignup}
                        sx={{ cursor: "pointer", textDecoration: "underline", fontSize: "1rem", color: "#1976d2" }}
                    >
                        Don't have a wallet? Generate One
                    </Link>
                </Typography>
            </Container>
        </>
    );
};

export default LoginPage;
