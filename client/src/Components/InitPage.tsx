// src/Components/InitPage.tsx

import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
    Container,
    Typography,
    CircularProgress,
    Button,
    Box,
    Alert,
} from "@mui/material";

const InitPage: React.FC = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    const initializeBtcService = async () => {
        setLoading(true);
        setError(null);

        try {
            const response = await fetch(`http://localhost:8080/api/btc/init`, { // Use relative path
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || "Initialization failed");
            }

            const data = await response.json();
            console.log("Initialization success:", data.message);

            // Save initialization status in localStorage
            localStorage.setItem("btc_initialized", "true");
            navigate("/signup");

        } catch (err: any) {
            console.error("Error initializing BTC Service:", err);
            setError(err.message || "An unexpected error occurred");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        const checkInitialization = async () => {
            const isInitialized = localStorage.getItem("btc_initialized");
            

            // if (isInitialized === "true") {
            //     navigate("/signup");
            //     return;
            // }

            await initializeBtcService();
        };

        checkInitialization();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return (
        <Container
            maxWidth="sm"
            style={{
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                justifyContent: "center",
                height: "100vh",
            }}
        >
            <Typography variant="h4" gutterBottom>
                Initializing Squid coin Service...
            </Typography>
            {loading && <CircularProgress />}
            {error && (
                <Box mt={2}>
                    <Alert severity="error">{error}</Alert>
                    <Box mt={2}>
                        <Button variant="contained" color="primary" onClick={initializeBtcService}>
                            Retry
                        </Button>
                    </Box>
                </Box>
            )}
        </Container>
    );
};

export default InitPage;
