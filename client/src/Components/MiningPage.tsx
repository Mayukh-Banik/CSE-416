import React, { useState, useEffect, useRef, useCallback } from "react";
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
    TextField,
} from "@mui/material";
import { styled, useTheme } from "@mui/material/styles";
import Sidebar from "./Sidebar";
import { useNavigate } from "react-router-dom";

const refreshInterval = 3000; // 3 seconds
const pauseDuration = 5000; // 5 seconds

// Styled components
const MiningContainer = styled(Box)(({ theme }) => ({
    padding: theme.spacing(4),
    transition: "margin-left 0.3s ease",
    [theme.breakpoints.down("sm")]: {
        marginLeft: 0,
    },
}));

interface MiningInfo {
    blocks: number;
    currentblocksize: number;
    currentblockweight: number;
    currentblocktx: number;
    difficulty: number;
    errors: string;
    generate: boolean;
}

interface MiningDashboardData {
    balance: string;
    miningInfo: MiningInfo;
}

interface ApiResponse<T> {
    status: string;
    message?: string;
    data: T;
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const MiningPage: React.FC = () => {
    const theme = useTheme();
    const navigate = useNavigate();
    const [isMining, setIsMining] = useState<boolean>(false);
    const [balance, setBalance] = useState<string>("initializing...");
    const [miningInfo, setMiningInfo] = useState<MiningInfo | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [numBlocks, setNumBlocks] = useState<number>(1);
    // Additional: State for disabling the button when in start state
    const [disableStartButton, setDisableStartButton] = useState<boolean>(false);
    const [buttonLock, setButtonLock] = useState<boolean>(false);

    const isTransitioningRef = useRef<boolean>(false);

    const API_BASE_URL = "http://localhost:8080/api/btc";

    const intervalRef = useRef<NodeJS.Timeout | null>(null);
    const timeoutRef = useRef<NodeJS.Timeout | null>(null);

    const fetchMiningDashboard = useCallback(async () => {
        try {
            const response = await fetch(`${API_BASE_URL}/miningdashboard`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || "Failed to fetch mining status");
            }

            const apiResponse: ApiResponse<MiningDashboardData> = await response.json();

            if (apiResponse.status !== "success") {
                throw new Error(apiResponse.message || "Failed to fetch mining dashboard");
            }

            setBalance(apiResponse.data.balance);
            setMiningInfo(apiResponse.data.miningInfo);

            setIsMining((prevStatus) => {
                if (prevStatus !== apiResponse.data.miningInfo.generate) {
                    setIsLoading(false);
                }
                return apiResponse.data.miningInfo.generate;
            });
        } catch (err) {
            setIsLoading(true);
            console.log("Timer started");
            setTimeout(() => {
                setIsLoading(false);
                console.log("Timer ended");
            }, 10000);
        }
    }, [API_BASE_URL]);

    useEffect(() => {
        fetchMiningDashboard();
        intervalRef.current = setInterval(fetchMiningDashboard, refreshInterval);

        return () => {
            if (intervalRef.current) clearInterval(intervalRef.current);
            if (timeoutRef.current) clearTimeout(timeoutRef.current);
        };
    }, [fetchMiningDashboard]);

    // Log state changes
    useEffect(() => {
        console.log("isLoading changed: ", isLoading);
    }, [isLoading]);

    useEffect(() => {
        console.log("isMining changed: ", isMining);
    }, [isMining]);

    useEffect(() => {
        console.log("disableStartButton changed: ", disableStartButton);
    }, [disableStartButton]);

    useEffect(() => {
        if (!isMining) {
            setIsLoading(false);
            // Disable Start button for 5 seconds
            setDisableStartButton(true);

            timeoutRef.current = setTimeout(() => {
                setDisableStartButton(false);
            }, pauseDuration);
        }
    }, [isMining]);

    const handleStartMining = async () => {
        if (numBlocks <= 0) {
            setError("Number of blocks must be greater than 0");
            return;
        }

        setIsLoading(true);
        console.log("Here");
        setButtonLock(true);
        try {
            const response = await fetch(`${API_BASE_URL}/startmining`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ numBlock: numBlocks }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                console.error("Server Error:", errorData);
                throw new Error(errorData.message || "Failed to start mining");
            }

            const data = await response.json();
            setSuccess(data.message || "Mining started successfully");
        } catch (error) {
            console.error("Error starting mining:", error);
            setIsMining(false);
            setError(error instanceof Error ? error.message : "An unknown error occurred.");
        } finally {
            setTimeout(() => {
                console.log("Stop Mining: Finished");
                setIsLoading(false);
                setButtonLock(false);
            }, 10000);
        }
    };

    const handleStopMining = async () => {
        console.log("Stop Mining: Button clicked");
        setIsLoading(true);
        console.log("Here");
        setButtonLock(true);

        try {
            const response = await fetch(`${API_BASE_URL}/stopmining`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                console.error("Server Error:", errorData);
                throw new Error(errorData.message || "Failed to stop mining");
            }

            const data = await response.json();
            console.log("Mining stopped successfully:", data);
            setSuccess(data.message || "Mining stopped successfully");

            // Clear interval: pause fetch for 5 seconds
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
                intervalRef.current = null;
            }

            // Resume fetchMiningDashboard after 5 seconds
            timeoutRef.current = setTimeout(() => {
                fetchMiningDashboard();
                intervalRef.current = setInterval(fetchMiningDashboard, refreshInterval);
            }, pauseDuration);
        } catch (error) {
            console.error("Error stopping mining:", error);
            setError(error instanceof Error ? error.message : "An unknown error occurred.");
        } finally {
            setTimeout(() => {
                console.log("Stop Mining: Finished");
                setIsLoading(false);
                setButtonLock(false);
            }, 10000);
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
                Mining Dashboard
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
                                <Box
                                    sx={{
                                        display: "flex",
                                        alignItems: "center",
                                        gap: 2,
                                        marginTop: 2,
                                    }}
                                >
                                    <TextField
                                        type="number"
                                        label="Blocks to Mine"
                                        value={numBlocks}
                                        onChange={(e) => setNumBlocks(Number(e.target.value))}
                                        sx={{ width: 150 }}
                                        disabled={isMining || isLoading}
                                        inputProps={{ min: 1 }}
                                    />
                                    <Button
                                        variant="contained"
                                        color="primary"
                                        onClick={isMining ? handleStopMining : handleStartMining}
                                        disabled={buttonLock && (isLoading || (!isMining && disableStartButton))}
                                        sx={{ height: '56px' }}
                                    >
                                        {isLoading ? (
                                            <CircularProgress size={24} />
                                        ) : isMining ? (
                                            "Stop Mining"
                                        ) : (
                                            "Start Mining"
                                        )}
                                    </Button>
                                </Box>
                            </Box>
                        </Paper>
                    </Grid>

                    {/* Mining Status */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6">Current Mining Status</Typography>
                                <Typography variant="body1" color={isMining ? "green" : "textSecondary"}>
                                    {isMining ? "Mining in Progress" : "Idle"}
                                </Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Difficulty */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6">Difficulty</Typography>
                                <Typography variant="body1">{miningInfo?.difficulty || "N/A"}</Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Block Count */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6">Block Count</Typography>
                                <Typography variant="body1">{miningInfo?.blocks || "N/A"}</Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Current Block Size */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6">Current Block Size</Typography>
                                <Typography variant="body1">{miningInfo?.currentblocksize || "N/A"} Kilobyte</Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Current Block Weight */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6">Current Block Weight</Typography>
                                <Typography variant="body1">{miningInfo?.currentblockweight || "N/A"} WU</Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Current Block Transactions */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6">Current Block Transactions</Typography>
                                <Typography variant="body1">{miningInfo?.currentblocktx || "N/A"}</Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Balance */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    Available Balance
                                </Typography>
                                <Typography variant="body1">{balance} BTC</Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Errors */}
                    {miningInfo?.errors && miningInfo.errors.trim() !== "" && (
                        <Grid item xs={12}>
                            <Card>
                                <CardContent>
                                    <Typography variant="h6" gutterBottom color="error">
                                        Mining Errors
                                    </Typography>
                                    <Typography variant="body1">{miningInfo.errors}</Typography>
                                </CardContent>
                            </Card>
                        </Grid>
                    )}
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
