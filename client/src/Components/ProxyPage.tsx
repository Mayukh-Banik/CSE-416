import React, { useState, useEffect } from 'react';
import Sidebar from "./Sidebar";
import useProxyHostsStyles from '../Stylesheets/ProxyPageStyles';
import { Button, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography, Box, TextField, LinearProgress } from '@mui/material';

interface ProxyHost {
  name: string;
  location: string;
  logs: string[];
  Statistics: {
    uptime: string;
  };
  address: string,
  peer_id: string;
  bandwidth: string;
  isEnabled: boolean;
  price: string;
  isHost: boolean;
  WalletAddressToSend: string; // New field

}

function getPrivateIP(callback: (ip: string | null) => void) {
  const peerConnection = new RTCPeerConnection({ iceServers: [] });
  peerConnection.createDataChannel('');
  peerConnection.createOffer()
    .then(offer => peerConnection.setLocalDescription(offer))
    .catch(err => console.error('Error creating offer:', err));

  peerConnection.onicecandidate = (event) => {
    if (event.candidate) {
      const parts = event.candidate.candidate.split(' ');
      const ip = parts[4]; // Extracts the IP address
      callback(ip);
      peerConnection.close();
    } else {
      console.log('No candidate found.');
    }
  };

}


const ProxyHosts: React.FC = () => {
  const [data, setData] = useState<ProxyHost[]>([]); // Corrected to ProxyHost array
  const [input, setInput] = useState<string>('');
  const [proxyHosts, setProxyHosts] = useState<ProxyHost[]>([]); // State to store proxy hosts
  const [currentIP, setCurrentIP] = useState<string>('');
  const [connectedProxy, setConnectedProxy] = useState<ProxyHost | null>(null);
  const [proxyHistory, setProxyHistory] = useState<{ name: string; location: string; timestamp: string }[]>([]);
  const [showHistoryOnly, setShowHistoryOnly] = useState<boolean>(false);
  const [showForm, setShowForm] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false); // Track loading state
  const [isHosting, setIsHosting] = useState<boolean>(false); // Track hosting state
  const [input1, setInput1] = useState('');
  const [input2, setInput2] = useState('');
  const [input3, setInput3] = useState('');
  const [input4, setInput4] = useState('');

  const [newProxy, setNewProxy] = useState<ProxyHost>({
    name: '',
    location: '',
    logs: [],
    address: '',
    Statistics: { uptime: '' },
    bandwidth: '',
    peer_id: '',
    isEnabled: false,
    price: '',
    isHost: false,
    WalletAddressToSend: '',
  });
  const styles = useProxyHostsStyles();

  const fetchData = async () => {
    try {
      setLoading(true);
      const response = await fetch('http://localhost:8081/proxy-data/', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const rawText = await response.text();
      const cleanedText = rawText.replace(/^null\s*/, '').trim();
      const result = JSON.parse(cleanedText);
      console.log("Parsed reult", result)
      const proxyData = Array.isArray(result) ? result : [result];
      const isEmptyProxy = (proxy: ProxyHost) => {
        return !proxy.name && !proxy.location && !proxy.price && !proxy.Statistics.uptime && !proxy.bandwidth;
      };
      const nonEmptyProxies = proxyData.filter(proxy => !isEmptyProxy(proxy));
      const hostproxy = proxyData.filter(proxy => !isEmptyProxy(proxy) && proxy.isHost);
      const ip = hostproxy.map(proxy => proxy.address);
      // Extract the address of the proxy where isHost is true
      setCurrentIP(ip[0]);
      setProxyHosts(nonEmptyProxies.map(proxy => ({
        name: proxy.name,
        location: proxy.location,
        address: proxy.address,
        price: proxy.price,
        Statistics: { uptime: proxy.Statistics?.uptime },
        bandwidth: proxy.bandwidth,
        logs: proxy.logs,
        isEnabled: false,
        isHost: proxy.isHost || false,
        peer_id: proxy.peer_id || '',
        WalletAddressToSend: proxy.WalletAddressToSend || ''
      })));
      getPrivateIP((privateIP) => {
        if (privateIP) {
          setCurrentIP(privateIP); // Set the private IP if found
        } else {
          console.error('Failed to retrieve private IP.');
        }
      });
    } catch (error) {
      console.error('Error fetching or processing data:', error);
    }
    finally {
      setLoading(false); // Stop loading
    }
  };
  const sendData = async () => {
    try {
      const response = await fetch('http://localhost:8081/proxy-data/', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newProxy),
      });

      if (!response.ok) {
        throw new Error('Failed to add proxy');
      }

      const result = await response.json();

      console.log('Proxy added successfully:', result);
    } catch (error) {
      console.error('Error adding proxy:', error);
      alert('There was an error adding the proxy. Please try again.');
    }
  };

  useEffect(() => {
    getPrivateIP((ip) => {
      if (ip) {
        setCurrentIP(ip);
      } else {
        console.error('Unable to retrieve private IP');
      }
    });


    fetchData()
  }, []);

  const handleDisconnect = (host: ProxyHost) => {
    setConnectedProxy(host);
    notifyDisConnectionToBackend(host);
    const newHistoryEntry = { name: host.name, location: host.location, timestamp: new Date().toLocaleString() };
    setProxyHistory([...proxyHistory, newHistoryEntry]);
  }

  const notifyDisConnectionToBackend = async (host: ProxyHost) => {
    console.log("Attempting to disconnect...");
    try {
      console.log(host.peer_id)
      console.log(host.address)
      const response = await fetch(`http://localhost:8081/disconnect-from-proxy?val=${host.peer_id}&ip=${host.address}`, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          hostName: host.name,
          hostLocation: host.location,
          hostPeerID: host.peer_id,
          hostAddress: host.address,

          timestamp: new Date().toLocaleString(),
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to notify backend about the disconnection');
      }
      alert(`Connected to ${host.location}`);
      if (response.ok) {
        setCurrentIP(host.address);
      }
      console.log(`Successfully notified backend about the disconnection to ${host.location}`);
    } catch (error) {
      window.alert("ERROR DISCONNECTING TO PROXY")
      console.error('Error notifying backend:', error);
    }
  };

  const handleConnect = (host: ProxyHost) => {
    console.log('Input 1:', input1);
    console.log('Input 2:', input2);
    const updatedHosts = proxyHosts.map(h => ({
      ...h,
      isEnabled: h.location === host.location ? true : h.isEnabled,
    }));

    setProxyHosts(updatedHosts);
    setConnectedProxy(host);
    notifyConnectionToBackend(host);


    // Update proxy history
    const newHistoryEntry = { name: host.name, location: host.location, timestamp: new Date().toLocaleString() };
    setProxyHistory([...proxyHistory, newHistoryEntry]);

  }

  const notifyConnectionToBackend = async (host: ProxyHost) => {
    console.log("Attempting to connect...");
    try {
      console.log(host.address)
      console.log(host.address)
      const response = await fetch(`http://localhost:8081/connect-proxy/`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          hostName: host.name,
          hostLocation: host.location,
          hostPeerID: host.peer_id,
          proxyIP: host.address,
          timestamp: new Date().toLocaleString(),
          passphrase: input1,
          transactionID: input2,
          destinationAddress: host.WalletAddressToSend,
          amount: input4
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to notify backend about the connection');
      }
      alert(`Connected to ${host.location}`);
      if (response.ok) {
        setCurrentIP(host.address);
      } if (response.ok) {
        setCurrentIP(host.address);
      }
      console.log(`Successfully notified backend about the connection to ${host.location}`);
    } catch (error) {
      window.alert("ERROR CONNECTING TO PROXY")
      console.error('Error notifying backend:', error);
    }
  };

  const handleAddProxy = async () => {
    newProxy.location = 'nyc'
    newProxy.Statistics.uptime = '0'
    newProxy.bandwidth = '0'
    if (newProxy.location.trim() === '' || newProxy.price.trim() === '' || newProxy.Statistics.uptime === '' || newProxy.bandwidth.trim() === '') {
      alert('Please fill in all fields.');
      return;
    }
    await sendData();

    setProxyHosts([...proxyHosts, { ...newProxy, logs: [], isEnabled: false }]);

    setNewProxy({
      name: '',
      location: '',
      logs: [],
      address: '',
      peer_id: '',
      Statistics: { uptime: '' },
      bandwidth: '',
      isEnabled: false,
      price: '',
      isHost: false,
      WalletAddressToSend: '',
    });

    setShowForm(false);
    alert('Proxy added successfully!');
    window.location.reload();
  };

  const handleSortByLocation = () => {
    const sortedHosts = [...proxyHosts].sort((a, b) => a.location.localeCompare(b.location));
    setProxyHosts(sortedHosts);
  };

  const handleSortByPrice = () => {
    const sortedHosts = [...proxyHosts].sort((a, b) => {
      const priceA = a.price === 'Free' ? 0 : parseFloat(a.price.replace('$', ''));
      const priceB = b.price === 'Free' ? 0 : parseFloat(b.price.replace('$', ''));
      return priceA - priceB;
    });
    setProxyHosts(sortedHosts);
  };

  const fetchHistory = async () => {
    try {
      const response = await fetch('http://localhost:8081/proxy-history/', {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
      });

      console.log(response.json());

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const history = await response.json();
      setProxyHistory(history); // Assuming setProxyHistory is defined elsewhere
    } catch (error) {
      console.error("Failed to fetch proxy history:", error);
    }
  };
  const handleClearAndShowHistory = () => {
    setShowHistoryOnly(true); // Show history
    fetchHistory();

  };

  const handleReturn = () => {
    setShowHistoryOnly(false); // Return
  };

  const handleStopHosting = async () => {
    try {
      const response = await fetch('http://localhost:8081/stop-hosting', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });

      if (!response.ok) {
        throw new Error('Failed to stop hosting');
      }

      alert('Stopped hosting successfully');
      await fetchData(); // Refresh data to update UI
    } catch (error) {
      console.error('Error stopping hosting:', error);
      alert('Error stopping hosting. Please try again.');
    }
  };

  return (
    <div className={styles.container}>
      <Box className={styles.boxContainer}>
        <Sidebar />
        <Box sx={{ marginTop: 2 }}>
          <Typography variant="h4">Proxy</Typography>
          {loading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '200px' }}>
              <LinearProgress sx={{ width: '100%' }} /> {/* Progress bar when loading is true */}
            </Box>
          ) : (
            <>
              <Box className={styles.header}>
                <Typography variant="h6">Your Current IP: {currentIP}</Typography>
                {!showHistoryOnly && (
                  <Button variant="outlined" onClick={handleClearAndShowHistory}>
                    Show Proxy History
                  </Button>
                )}
              </Box>
              <br />
              {/* Show the history */}
              {showHistoryOnly ? (
                <>
                  <Box className={styles.historyContainer}>
                    <Typography variant="h5">Proxy Connection History</Typography>
                    <TableContainer className={styles.historyTable}>
                      <Table>
                        <TableHead>
                          <TableRow>
                            <TableCell>Name</TableCell>
                            {/* <TableCell>Location</TableCell> */}
                            <TableCell>Connected At</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {proxyHistory.map((entry, index) => (
                            <TableRow key={index}>
                              <TableCell>{entry.name}</TableCell>
                              {/* <TableCell>{entry.location}</TableCell> */}
                              <TableCell>{entry.timestamp}</TableCell>
                              <TableCell>
                                <Button
                                  variant="contained"
                                  onClick={() => {
                                    const host = proxyHosts.find(
                                      (h) => h.location === entry.location
                                    );
                                    if (host) handleConnect(host);
                                  }}
                                >
                                  Connect
                                </Button>
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </Box>
                  <Box sx={{ marginTop: 2 }}>
                    <Button variant="contained" onClick={handleReturn}>
                      Return to Proxy Hosts
                    </Button>
                  </Box>
                </>
              ) : (
                <>
                  <Box className={styles.proxyButton}>
                    {/* Add Yourself as Proxy */}
                    <Button
                      variant="contained"
                      onClick={() => setShowForm((prev) => !prev)}
                    >
                      {showForm ? 'Hide Form' : 'Add Yourself as Proxy'}
                    </Button>
                    <Button variant="outlined" color="error" onClick={handleStopHosting}>
                      Stop Hosting
                    </Button>
                    {/* Sort Buttons */}
                    <Box sx={{ display: 'flex', gap: '10px' }}>
                      {/* <Button variant="outlined" onClick={handleSortByLocation}>
                        Sort by Location
                      </Button> */}
                      <Button variant="outlined" onClick={handleSortByPrice}>
                        Sort by Price
                      </Button>
                    </Box>
                  </Box>
                  {/* Expandable Form Section */}
                  {showForm && (
                    <Box sx={{ marginTop: 1 }} className={styles.form}>
                      <Typography variant="h6">
                        Fill in your proxy details
                      </Typography>
                      <Box sx={{ display: 'flex', gap: 2 }}>
                        <TextField
                          label="Name"
                          variant="outlined"
                          value={newProxy.name}
                          onChange={(e) =>
                            setNewProxy({ ...newProxy, name: e.target.value })
                          }
                        />
                        {/* <TextField
                          label="Location"
                          variant="outlined"
                          value={newProxy.location}
                          onChange={(e) =>
                            setNewProxy({ ...newProxy, location: e.target.value })
                          }
                        /> */}
                        <TextField
                          label="Price ($)"
                          variant="outlined"
                          type="number"
                          value={newProxy.price}
                          onChange={(e) =>
                            setNewProxy({ ...newProxy, price: e.target.value })
                          }
                          InputProps={{ inputProps: { min: 0 } }}
                        />
                        {/* <TextField
                          label="Uptime (%)"
                          variant="outlined"
                          value={newProxy.Statistics.uptime}
                          onChange={(e) =>
                            setNewProxy({
                              ...newProxy,
                              Statistics: {
                                ...newProxy.Statistics,
                                uptime: e.target.value,
                              },
                            })
                          }
                        /> */}
                        {/* <TextField
                          label="Bandwidth"
                          variant="outlined"
                          value={newProxy.bandwidth}
                          onChange={(e) =>
                            setNewProxy({ ...newProxy, bandwidth: e.target.value })
                          }
                        /> */}
                        <Button
                          variant="contained"
                          className={styles.submitButton}
                          onClick={handleAddProxy}
                        >
                          Add Proxy
                        </Button>
                      </Box>
                    </Box>
                  )}
                  <TableContainer className={styles.proxyTable}>
                    <Table className={styles.table}>
                      <TableHead>
                        <TableRow>
                          <TableCell>Name</TableCell>
                          {/* <TableCell>Location</TableCell> */}
                          <TableCell>Price</TableCell>
                          {/* <TableCell>Uptime</TableCell> */}
                          {/* <TableCell>Bandwidth</TableCell> */}
                          {/* <TableCell>Logs</TableCell> */}
                          <TableCell>Actions</TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {proxyHosts.map((host, index) => (
                          <TableRow key={index}>
                            <TableCell>{host.name}</TableCell>
                            {/* <TableCell>{host.location}</TableCell> */}
                            <TableCell>{host.price}</TableCell>
                            {/* <TableCell>
                              {host.Statistics && host.Statistics.uptime
                                ? host.Statistics.uptime
                                : 'N/A'}
                            </TableCell> */}
                            {/* <TableCell>{host.bandwidth}</TableCell> */}
                            {/* <TableCell>
                              <Button
                                variant="text"
                                onClick={() => alert(host.logs.join('\n'))}
                              >
                                View Logs
                              </Button>
                            </TableCell> */}
                            <TableCell>
                              <TextField
                                label="Passphrase"
                                variant="outlined"
                                value={input1}
                                onChange={(e) => setInput1(e.target.value)}
                              />
                              <TextField
                                label="Transaction ID"
                                variant="outlined"
                                value={input2}
                                onChange={(e) => setInput2(e.target.value)}
                              />
                              <TextField
                                label="Amount"
                                variant="outlined"
                                type="number"
                                value={input4}
                                onChange={(e) => {
                                  const value = e.target.value;
                                  if (value === '' || parseFloat(value) >= 0) {
                                    setInput4(value);
                                  }
                                }}
                                inputProps={{
                                  min: 0,
                                  step: "any"
                                }}
                              />
                              <Button
                                variant="contained"
                                onClick={() => handleConnect(host)}
                              >
                                Connect
                              </Button>
                            </TableCell>
                            <TableCell>
                              <Button
                                variant="contained"
                                onClick={() => handleDisconnect(host)}
                              >
                                Disconnet
                              </Button>
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </TableContainer>
                </>
              )}
            </>
          )}
        </Box>
      </Box>
    </div>
  );
}
export default ProxyHosts;