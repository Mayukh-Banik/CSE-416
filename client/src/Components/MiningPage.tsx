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
    TextField,
} from "@mui/material";
import { styled, useTheme } from "@mui/material/styles";
import Sidebar from "./Sidebar";
import { useNavigate } from "react-router-dom";

const refreshInterval = 5000; // 5 seconds

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
    mining: boolean;
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const MiningPage: React.FC = () => {
    const theme = useTheme();
    const navigate = useNavigate();
    const [isMining, setIsMining] = useState<boolean>(false);

    const [balance, setBalance] = useState<string>("0 SQC");
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [numBlocks, setNumBlocks] = useState<number>(1); // 블록 수 상태

    const API_BASE_URL = "http://localhost:8080/api/btc"; // 백엔드 API 베이스 URL


    useEffect(() => {
        console.log("Mining status changed:", isMining);
        if (!isMining) {
            setIsLoading(false); // 마이닝 중단 시 로딩 해제
        }
    }, [isMining]);


    // Fetch mining status periodically
    useEffect(() => {
        const fetchMiningStatus = async () => {
            try {
                const response = await fetch(`${API_BASE_URL}/getminingstatus`, {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                    },
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.message || "Failed to fetch mining status");
                }

                const data = await response.json();

                // 이전 상태와 비교하여 변경된 경우만 처리
                setIsMining((prevStatus) => {
                    if (prevStatus !== data.data.mining) {
                        console.log("Mining status changed:", data.data.mining);
                        setIsLoading(false); // 로딩 상태 해제
                    }
                    return data.data.mining; // 새로운 상태 설정
                });
            } catch (err) {
                console.error("Error fetching mining status:", err);
                setError(err instanceof Error ? err.message : "An unknown error occurred.");
                setIsLoading(false); // 로딩 상태 해제
            }
        };

        fetchMiningStatus(); // 컴포넌트 초기 렌더링 시 실행
        const interval = setInterval(fetchMiningStatus, refreshInterval); // 5초마다 상태 갱신

        return () => clearInterval(interval); // 컴포넌트 언마운트 시 인터벌 해제
    }, []);

    // Fetch balance periodically
    useEffect(() => {
        const fetchBalance = async () => {
            try {
                const response = await fetch(`${API_BASE_URL}/balance`, {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                    },
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.message || "Failed to fetch balance");
                }

                const data = await response.json();
                setBalance(data.data.balance || "0 SQC");
            } catch (err) {
                console.error("Error fetching balance:", err);
                setError(err instanceof Error ? err.message : "An unknown error occurred.");
            }
        };

        fetchBalance(); // 컴포넌트가 처음 렌더링될 때 즉시 호출

        const interval = setInterval(fetchBalance, refreshInterval); // 5초마다 갱신
        return () => clearInterval(interval); // 컴포넌트 언마운트 시 인터벌 정리
    }, []);

    // Start mining
    const handleStartMining = async () => {
        if (numBlocks <= 0) {
            setError("Number of blocks must be greater than 0");
            return;
        }

        console.log("Start Mining: Button clicked"); // 디버깅 로그
        setIsLoading(true);

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
                console.error("Server Error:", errorData); // 서버 에러 로그
                throw new Error(errorData.message || "Failed to start mining");
            }

            const data = await response.json();
            console.log("Mining started successfully:", data); // 서버 성공 로그
            setSuccess(data.message || "Mining started successfully");
        } catch (error) {
            console.error("Error starting mining:", error); // 오류 로그
            setIsMining(false); // 상태 롤백
            setError(error instanceof Error ? error.message : "An unknown error occurred.");
        } finally {
            console.log("Start Mining: Finished"); // 디버깅 로그
            setIsLoading(false); // 로딩 상태 해제
        }
    };

    // Stop mining
    const handleStopMining = async () => {
        console.log("Stop Mining: Button clicked"); // 디버깅 로그
        setIsLoading(true);

        try {

            const response = await fetch(`${API_BASE_URL}/stopmining`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                console.error("Server Error:", errorData); // 서버 에러 로그
                throw new Error(errorData.message || "Failed to stop mining");
            }

            const data = await response.json();
            console.log("Mining stopped successfully:", data); // 서버 성공 로그
            setSuccess(data.message || "Mining stopped successfully");
        } catch (error) {
            console.error("Error stopping mining:", error); // 오류 로그
            setError(error instanceof Error ? error.message : "An unknown error occurred.");
        } finally {
            console.log("Stop Mining: Finished"); // 디버깅 로그
            setIsLoading(false); // 로딩 상태 해제
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
                                        disabled={isMining || isLoading} // 마이닝 상태는 비활성화하지 않음
                                        inputProps={{ min: 1 }}
                                    />
                                    <Button
                                        variant="contained"
                                        color="primary"
                                        onClick={isMining ? handleStopMining : handleStartMining}
                                        disabled={isLoading}
                                        sx={{ height: '56px' }} // TextField와 높이 맞춤
                                    >
                                        {isLoading ? (<CircularProgress size={24} />) : isMining ? ("Stop Mining") : ("Start Mining")}
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
                                <Typography variant="body1">{isMining ? "Mining in Progress" : "Idle"}</Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Mined Blocks - Placeholder */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    Mined Blocks
                                </Typography>
                                <Typography variant="body1" color="textSecondary">
                                    Coming Soon...
                                </Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Balance */}
                    <Grid item xs={12}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    Balance
                                </Typography>
                                <Typography variant="body1">{balance}</Typography>
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
