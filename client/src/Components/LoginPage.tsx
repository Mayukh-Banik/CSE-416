import React, { useState, useEffect } from 'react';
import { TextField, Button, Container, Typography, Link, Box } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import useRegisterPageStyles from '../Stylesheets/RegisterPageStyles';
import Header from './Header';
import { useNavigate } from 'react-router-dom';
import forge from 'node-forge';

const PORT = 8080;

const LoginPage: React.FC = () => {
    const classes = useRegisterPageStyles();
    const navigate = useNavigate();

    const [publicKey, setPublicKey] = useState('');
    const [privateKey, setPrivateKey] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [challenge, setChallenge] = useState<string | null>(null);

    // 페이지 로드 시 챌린지 자동 요청
    useEffect(() => {
        const requestChallenge = async () => {
            try {
                const response = await fetch(`http://localhost:${PORT}/api/auth/request-challenge`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                if (!response.ok) {
                    const errorMessage = await response.text();
                    console.error("Server Error:", errorMessage);
                    throw new Error('Failed to request challenge');
                }

                const data = await response.json();
                setChallenge(data.challenge);
                console.log("Challenge received:", data.challenge);
            } catch (err) {
                console.error("Error requesting challenge:", err);
                setError('Failed to request challenge.');
            }
        };

        requestChallenge();
    }, []);

    // 서명 및 인증 함수
    const handleLogin = async () => {
        if (!publicKey || !privateKey || !challenge) {
            setError('Public key, private key, and challenge are required.');
            return;
        }

        try {
            // 프라이빗 키 파싱
            const rsaPrivateKey = forge.pki.privateKeyFromPem(privateKey);
            console.log('Private key parsed successfully:', rsaPrivateKey);

            // 챌린지 디코딩
            const challengeBytes = forge.util.decode64(challenge);
            console.log('Decoded challenge bytes:', challengeBytes);

            // 서명 생성
            const md = forge.md.sha256.create();
            md.update(challengeBytes);
            const signature = rsaPrivateKey.sign(md);
            console.log('Signature generated:', signature);

            // 서명을 Base64로 인코딩
            const signatureBase64 = forge.util.encode64(signature);

            // 서버로 서명과 퍼블릭 키 전송
            const response = await fetch(`http://localhost:${PORT}/api/auth/verify-challenge`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ public_key: publicKey, signature: signatureBase64 })
            });

            if (!response.ok) {
                const errorMessage = await response.text();
                console.error('Server Error:', errorMessage);
                throw new Error('Failed to verify signature');
            }

            const data = await response.json();
            console.log('Verification Response:', data);
            // 로그인 성공 처리 (예: 토큰 저장, 페이지 이동 등)
        } catch (err) {
            console.error('Error during login:', err);
            setError('Failed to login.');
        }
    };

    const handleSignup = () => {
        navigate('/');
    };

    return (
        <>
            <Header />
            <Container maxWidth="sm" sx={{ mt: 4 }}>
                <LockOutlinedIcon sx={{ fontSize: 40, mb: 2 }} />
                <Typography variant="h4" component="h1" gutterBottom>
                    Login
                </Typography>
                {error && <Typography color="error">{error}</Typography>}
                <Box sx={{ mt: 2 }}>
                    <TextField
                        label="Public Key"
                        multiline
                        rows={6}
                        fullWidth
                        value={publicKey}
                        onChange={(e) => setPublicKey(e.target.value)}
                        sx={{ mb: 2 }}
                    />
                    <TextField
                        label="Private Key"
                        multiline
                        rows={6}
                        fullWidth
                        value={privateKey}
                        onChange={(e) => setPrivateKey(e.target.value)}
                    />
                    <Button variant="contained" fullWidth sx={{ mt: 2 }} onClick={handleLogin}>
                        Login
                    </Button>
                </Box>
                <Typography sx={{ mt: 2 }}>
                    Don't have a wallet?{' '}
                    <Link onClick={handleSignup} sx={{ cursor: 'pointer' }}>
                        Generate One
                    </Link>
                </Typography>
            </Container>
        </>
    );
};

export default LoginPage;
