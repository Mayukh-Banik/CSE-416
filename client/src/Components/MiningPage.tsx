import React, { useState, useEffect } from "react";
import {
    Box,
    Typography,
    Button,
    CircularProgress,
    Snackbar,
    Alert,
    Paper,
    Grid,
    Card,
    CardContent,
} from "@mui/material";
import { styled, useTheme } from "@mui/material/styles";
import Sidebar from "./Sidebar";

// Styled components
const MiningContainer = styled(Box)(({ theme }) => ({
    padding: theme.spacing(4),
    transition: "margin-left 0.3s ease",
    [theme.breakpoints.down("sm")]: {
        marginLeft: 0,
    },
}));

// Define MiningStatus interface
interface MiningStatus {
    minedBlocks: number;
    lastMinedBlock: string;
    isMining: boolean;
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const MiningPage: React.FC = () => {
    const theme = useTheme();
    const [isMining, setIsMining] = useState(false);
    const [miningStatus, setMiningStatus] = useState<MiningStatus>({
        minedBlocks: 0,
        lastMinedBlock: "",
        isMining: false,
    });
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const API_BASE_URL = "http://localhost:8081"; 

    // Fetch mining status periodically
    useEffect(() => {
        const fetchMiningStatus = async () => {
            try {
                const response = await fetch(`${API_BASE_URL}/api/mining-status`);
                if (!response.ok) throw new Error("Failed to fetch mining status");
                const data: MiningStatus = await response.json();
                setMiningStatus(data);
                setIsMining(data.isMining);
            } catch (err) {
                setError(err instanceof Error ? err.message : "An unknown error occurred.");
            }
        };

        fetchMiningStatus();
        const interval = setInterval(fetchMiningStatus, 10000); // Set to 10 seconds for efficiency
        return () => clearInterval(interval);
    }, []);

    // Start mining
    const handleStartMining = async () => {
        try {
            const response = await fetch(`${API_BASE_URL}/api/start-mining`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
            });
            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || "Failed to start mining");
            }
            const data: MiningStatus = await response.json();
            setSuccess("Mining started");
            setMiningStatus(data);
            setIsMining(data.isMining);
        } catch (error) {
            setError(error instanceof Error ? error.message : "An unknown error occurred.");
        }
    };

    // Stop mining
    const handleStopMining = async () => {
        try {
            const response = await fetch(`${API_BASE_URL}/api/stop-mining`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
            });
            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || "Failed to stop mining");
            }
            const data: MiningStatus = await response.json();
            setSuccess("Mining stopped");
            setMiningStatus(data);
            setIsMining(data.isMining);
        } catch (error) {
            setError(error instanceof Error ? error.message : "An unknown error occurred.");
        }
    };

    const handleCloseSnackbar = () => {
        setError(null);
        setSuccess(null);
    };

    return (
        <Box
            sx={{
                padding: 2,
                marginTop: '70px',
                marginLeft: `${drawerWidth}px`,
                transition: 'margin-left 0.3s ease',
                [theme.breakpoints.down('sm')]: {
                    marginLeft: `${collapsedDrawerWidth}px`,
                },
            }}
        >
            <Sidebar />
            <Typography variant="h4" gutterBottom>
                Mining
            </Typography>
            <MiningContainer>
                <Grid container spacing={4}>
                    <Grid item xs={12}>
                        <Paper elevation={3} sx={{ padding: 4 }}>
                            <Box
                                sx={{
                                    display: "flex",
                                    flexDirection: "column",
                                    alignItems: "center",
                                }}
                            >
                                <Typography variant="h6" gutterBottom>
                                    ~ SquidCoins ~
                                </Typography>
                                <Button
                                    variant="contained"
                                    color="primary"
                                    onClick={isMining ? handleStopMining : handleStartMining}
                                    sx={{ mb: 2 }}
                                >
                                    {isMining ? "Stop Mining" : "Start Mining"}
                                </Button>
                                {isMining && <CircularProgress />}
                            </Box>
                        </Paper>
                    </Grid>

                    {/* Mining Status */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    Current Mining Status
                                </Typography>
                                <Typography variant="body1">
                                    {isMining ? "Mining in Progress" : "Idle"}
                                </Typography>
                                <Typography variant="body2" color="textSecondary">
                                    Last Mined Block: {miningStatus.lastMinedBlock || "N/A"}
                                </Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Mined Blocks */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    Mined Blocks
                                </Typography>
                                <Typography variant="h4">{miningStatus.minedBlocks}</Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Additional Information or Charts */}
                    <Grid item xs={12}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    Balance
                                </Typography>
                                <Typography variant="body1">
                                    100 SQC
                                </Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                </Grid>

                {/* Feedback */}
                <Snackbar
                    open={!!error}
                    autoHideDuration={6000}
                    onClose={handleCloseSnackbar}
                    anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
                >
                    <Alert onClose={handleCloseSnackbar} severity="error" sx={{ width: "100%" }}>
                        {error}
                    </Alert>
                </Snackbar>

                <Snackbar
                    open={!!success}
                    autoHideDuration={6000}
                    onClose={handleCloseSnackbar}
                    anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
                >
                    <Alert onClose={handleCloseSnackbar} severity="success" sx={{ width: "100%" }}>
                        {success}
                    </Alert>
                </Snackbar>
            </MiningContainer>
        </Box>
    );
};

export default MiningPage;
