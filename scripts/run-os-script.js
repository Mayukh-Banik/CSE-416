const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

const scripts = {
    buildBtcd: {
        windows: () => {
            const btcdDir = path.resolve("application-layer", "btcd");
            console.log(`Changing directory to: ${btcdDir}`);
            try {
                process.chdir(btcdDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go build -o btcd.exe";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("BTCD built successfully.");
            } catch (err) {
                console.error(`Failed to build btcd.exe: ${err.message}`);
                process.exit(1);
            }
        },
        macos: () => {
            const btcdDir = path.resolve("application-layer", "btcd");
            console.log(`Changing directory to: ${btcdDir}`);
            try {
                process.chdir(btcdDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go build -o btcd";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("BTCD built successfully.");
            } catch (err) {
                console.error(`Failed to build btcd: ${err.message}`);
                process.exit(1);
            }
        },
        linux: () => {
            const btcdDir = path.resolve("application-layer", "btcd");
            console.log(`Changing directory to: ${btcdDir}`);
            try {
                process.chdir(btcdDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go build -o btcd";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("BTCD built successfully.");
            } catch (err) {
                console.error(`Failed to build btcd: ${err.message}`);
                process.exit(1);
            }
        }
    },
    startBtcd: {
        windows: () => {
            const btcdDataPath = path.join(
                process.env.USERPROFILE,
                "AppData",
                "Local",
                "Btcd",
                "data",
                "mainnet",
                "blocks_ffldb"
            );
            console.log(`Deleting corrupted data at: ${btcdDataPath}`);
            try {
                fs.rmSync(btcdDataPath, { recursive: true, force: true });
                console.log("Successfully deleted corrupted data.");
            } catch (err) {
                console.error(`Failed to delete corrupted data: ${err.message}`);
                process.exit(1);
            }

            const btcdDir = path.resolve("application-layer", "btcd");
            const btcdExe = path.join(btcdDir, "btcd.exe");
            const configFile = path.join(btcdDir, "btcd.conf");

            console.log(`Changing directory to: ${btcdDir}`);
            try {
                process.chdir(btcdDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            console.log(`Executing: ${btcdExe} --configfile=${configFile}`);
            try {
                execSync(`"${btcdExe}" --configfile="${configFile}"`, { stdio: "inherit" });
                console.log("BTCD started successfully.");
            } catch (err) {
                console.error(`Failed to execute btcd.exe: ${err.message}`);
                process.exit(1);
            }
        },
        macos: () => {
            const btcdDataPath = path.join(
                os.homedir(),
                "Library",
                "Application Support",
                "Btcd",
                "data",
                "mainnet",
                "blocks_ffldb"
            );
            console.log(`Deleting corrupted data at: ${btcdDataPath}`);
            try {
                execSync(`rm -rf "${btcdDataPath}"`, { stdio: "inherit" });
                console.log("Successfully deleted corrupted data.");
            } catch (err) {
                console.error(`Failed to delete corrupted data: ${err.message}`);
                process.exit(1);
            }

            const btcdDir = path.resolve("application-layer", "btcd");
            const btcdExe = path.join(btcdDir, "btcd");
            const configFile = path.join(btcdDir, "btcd.conf");

            console.log(`Changing directory to: ${btcdDir}`);
            try {
                process.chdir(btcdDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            console.log(`Executing: ${btcdExe} --configfile=${configFile}`);
            try {
                execSync(`"${btcdExe}" --configfile="${configFile}"`, { stdio: "inherit" });
                console.log("BTCD started successfully.");
            } catch (err) {
                console.error(`Failed to execute btcd: ${err.message}`);
                process.exit(1);
            }
        },
        linux: () => {
            const btcdDataPath = path.join(
                os.homedir(),
                ".btcd",
                "data",
                "mainnet",
                "blocks_ffldb"
            );
            console.log(`Deleting corrupted data at: ${btcdDataPath}`);
            try {
                execSync(`rm -rf "${btcdDataPath}"`, { stdio: "inherit" });
                console.log("Successfully deleted corrupted data.");
            } catch (err) {
                console.error(`Failed to delete corrupted data: ${err.message}`);
                process.exit(1);
            }

            const btcdDir = path.resolve("application-layer", "btcd");
            const btcdExe = path.join(btcdDir, "btcd");
            const configFile = path.join(btcdDir, "btcd.conf");

            console.log(`Changing directory to: ${btcdDir}`);
            try {
                process.chdir(btcdDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            console.log(`Executing: ${btcdExe} --configfile=${configFile}`);
            try {
                execSync(`"${btcdExe}" --configfile="${configFile}"`, { stdio: "inherit" });
                console.log("BTCD started successfully.");
            } catch (err) {
                console.error(`Failed to execute btcd: ${err.message}`);
                process.exit(1);
            }
        }
    },
    buildBtcwallet: {
        windows: () => {
            const walletDir = path.resolve("application-layer", "btcwallet");
            console.log(`Changing directory to: ${walletDir}`);
            try {
                process.chdir(walletDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go build -o btcwallet.exe";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("BTCWallet built successfully.");
            } catch (err) {
                console.error(`Failed to build btcwallet.exe: ${err.message}`);
                process.exit(1);
            }
        },
        macos: () => {
            const walletDir = path.resolve("application-layer", "btcwallet");
            console.log(`Changing directory to: ${walletDir}`);
            try {
                process.chdir(walletDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go build -o btcwallet";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("BTCWallet built successfully.");
            } catch (err) {
                console.error(`Failed to build btcwallet: ${err.message}`);
                process.exit(1);
            }
        },
        linux: () => {
            const walletDir = path.resolve("application-layer", "btcwallet");
            console.log(`Changing directory to: ${walletDir}`);
            try {
                process.chdir(walletDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go build -o btcwallet";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("BTCWallet built successfully.");
            } catch (err) {
                console.error(`Failed to build btcwallet: ${err.message}`);
                process.exit(1);
            }
        }
    },
    startBtcwallet: {
        windows: () => {
            const walletDir = path.resolve("application-layer", "btcwallet");
            const walletExe = path.join(walletDir, "btcwallet.exe");
            const configFile = path.join(walletDir, "btcwallet.conf");

            console.log(`Changing directory to: ${walletDir}`);
            try {
                process.chdir(walletDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            console.log(`Executing: ${walletExe} --configfile=${configFile}`);
            try {
                execSync(`"${walletExe}" --configfile="${configFile}"`, { stdio: "inherit" });
                console.log("BTCWallet started successfully.");
            } catch (err) {
                console.error(`Failed to execute btcwallet.exe: ${err.message}`);
                process.exit(1);
            }
        },
        macos: () => {
            const walletDir = path.resolve("application-layer", "btcwallet");
            const walletExe = path.join(walletDir, "btcwallet");
            const configFile = path.join(walletDir, "btcwallet.conf");

            console.log(`Changing directory to: ${walletDir}`);
            try {
                process.chdir(walletDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            console.log(`Executing: ${walletExe} --configfile=${configFile}`);
            try {
                execSync(`"${walletExe}" --configfile="${configFile}"`, { stdio: "inherit" });
                console.log("BTCWallet started successfully.");
            } catch (err) {
                console.error(`Failed to execute btcwallet: ${err.message}`);
                process.exit(1);
            }
        },
        linux: () => {
            const walletDir = path.resolve("application-layer", "btcwallet");
            const walletExe = path.join(walletDir, "btcwallet");
            const configFile = path.join(walletDir, "btcwallet.conf");

            console.log(`Changing directory to: ${walletDir}`);
            try {
                process.chdir(walletDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            console.log(`Executing: ${walletExe} --configfile=${configFile}`);
            try {
                execSync(`"${walletExe}" --configfile="${configFile}"`, { stdio: "inherit" });
                console.log("BTCWallet started successfully.");
            } catch (err) {
                console.error(`Failed to execute btcwallet: ${err.message}`);
                process.exit(1);
            }
        }
    },
    buildBtcctl: {
        windows: () => {
            const ctlDir = path.resolve("application-layer", "btcd", "cmd", "btcctl");
            console.log(`Changing directory to: ${ctlDir}`);
            try {
                process.chdir(ctlDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go build -o btcctl.exe";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("BTCCTL built successfully.");
            } catch (err) {
                console.error(`Failed to build btcctl.exe: ${err.message}`);
                process.exit(1);
            }
        },
        macos: () => {
            const ctlDir = path.resolve("application-layer", "btcd", "cmd", "btcctl");
            console.log(`Changing directory to: ${ctlDir}`);
            try {
                process.chdir(ctlDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go build -o btcctl";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("BTCCTL built successfully.");
            } catch (err) {
                console.error(`Failed to build btcctl: ${err.message}`);
                process.exit(1);
            }
        },
        linux: () => {
            const ctlDir = path.resolve("application-layer", "btcd", "cmd", "btcctl");
            console.log(`Changing directory to: ${ctlDir}`);
            try {
                process.chdir(ctlDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go build -o btcctl";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("BTCCTL built successfully.");
            } catch (err) {
                console.error(`Failed to build btcctl: ${err.message}`);
                process.exit(1);
            }
        }
    },
    testWallet: {
        windows: () => {
            console.log("Running wallet tests on Windows...");
            const walletDir = path.resolve("application-layer");
            console.log(`Changing directory to: ${walletDir}`);
            try {
                process.chdir(walletDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go test ./wallet -v -count=1";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("Wallet tests ran successfully.");
            } catch (err) {
                console.error(`Failed to run wallet tests: ${err.message}`);
                process.exit(1);
            }
        },
        macos: () => {
            console.log("Running wallet tests on macOS...");
            const walletDir = path.resolve("application-layer");
            console.log(`Changing directory to: ${walletDir}`);
            try {
                process.chdir(walletDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go test ./wallet -v -count=1";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("Wallet tests ran successfully.");
            } catch (err) {
                console.error(`Failed to run wallet tests: ${err.message}`);
                process.exit(1);
            }
        },
        linux: () => {
            console.log("Running wallet tests on Linux...");
            const walletDir = path.resolve("application-layer");
            console.log(`Changing directory to: ${walletDir}`);
            try {
                process.chdir(walletDir);
                console.log(`Current directory: ${process.cwd()}`);
            } catch (err) {
                console.error(`Failed to change directory: ${err.message}`);
                process.exit(1);
            }

            const command = "go test ./wallet -v -count=1";
            console.log(`Executing: ${command}`);
            try {
                execSync(command, { stdio: "inherit" });
                console.log("Wallet tests ran successfully.");
            } catch (err) {
                console.error(`Failed to run wallet tests: ${err.message}`);
                process.exit(1);
            }
        }
    }
};


const scriptName = process.argv[2];
if (!scriptName || !scripts[scriptName]) {
    console.error("Invalid script name! Available scripts:", Object.keys(scripts).join(", "));
    process.exit(1);
}

const platform = os.platform();
let osType = "";
if (platform === "win32") osType = "windows";
else if (platform === "darwin") osType = "macos";
else if (platform === "linux") osType = "linux";
else {
    console.error(`Unsupported platform: ${platform}`);
    process.exit(1);
}

const scriptFunction = scripts[scriptName][osType];
if (!scriptFunction) {
    console.error(`No script found for '${scriptName}' on '${osType}'`);
    process.exit(1);
}

try {
    console.log(`Running script '${scriptName}' on '${osType}'`);
    scriptFunction();
} catch (err) {
    console.error(`Error executing script '${scriptName}': ${err.message}`);
    process.exit(1);
}
