import React, { useState, useEffect } from 'react';
import Sidebar from "./Sidebar";
import useProxyHostsStyles from '../Stylesheets/ProxyPageStyles';
import { Button, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography, Box, TextField } from '@mui/material';

interface ProxyHost {
  name: string;
  location: string;
  logs: string[];
  peer_id: string;
  isEnabled: boolean;
  price: string;
  isHost: boolean;
}

const ProxyHosts: React.FC = () => {
  const [data, setData] = useState<ProxyHost[]>([]);
  const [input, setInput] = useState<string>('');
  const [proxyHosts, setProxyHosts] = useState<ProxyHost[]>([]);
  const [currentIP, setCurrentIP] = useState<string>('');
  const [connectedProxy, setConnectedProxy] = useState<ProxyHost | null>(null);
  const [proxyHistory, setProxyHistory] = useState<{ name: string; location: string; timestamp: string }[]>([]);
  const [showHistoryOnly, setShowHistoryOnly] = useState<boolean>(false);
  const [showForm, setShowForm] = useState<boolean>(false);
  const [newProxy, setNewProxy] = useState<ProxyHost>({
    name: '',
    location: '',
    logs: [],
    peer_id: '',
    isEnabled: false,
    price: '',
    isHost: false,
  });
  const styles = useProxyHostsStyles();

  const fetchData = async () => {
    try {
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
      console.log("Parsed result", result);
      const proxyData = Array.isArray(result) ? result : [result];
      const isEmptyProxy = (proxy: ProxyHost) => {
        return !proxy.name && !proxy.location && !proxy.price;
      };
      const nonEmptyProxies = proxyData.filter(proxy => !isEmptyProxy(proxy));

      setProxyHosts(nonEmptyProxies.map(proxy => ({
        name: proxy.name,
        location: proxy.location || proxy.address,
        price: proxy.price,
        logs: proxy.logs,
        isEnabled: false,
        isHost: proxy.isHost || false,
        peer_id: proxy.peer_id || '',
      })));
    } catch (error) {
      console.error('Error fetching or processing data:', error);
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
    const fetchIP = async () => {
      try {
        const response = await fetch('https://api.ipify.org?format=json');
        const data = await response.json();
        setCurrentIP(data.ip);
      } catch (error) {
        console.error("Error fetching IP:", error);
      }
    };

    fetchIP();
    fetchData();
  }, []);

  const handleConnect = (host: ProxyHost) => {
    const updatedHosts = proxyHosts.map(h => ({
      ...h,
      isEnabled: h.location === host.location ? true : h.isEnabled,
    }));

    setProxyHosts(updatedHosts);
    setConnectedProxy(host);
    notifyConnectionToBackend(host);

    const newHistoryEntry = { name: host.name, location: host.location, timestamp: new Date().toLocaleString() };
    setProxyHistory([...proxyHistory, newHistoryEntry]);

    alert(`Connected to ${host.location}`);
  };

  const notifyConnectionToBackend = async (host: ProxyHost) => {
    console.log("Attempting to connect...");
    try {
      console.log(host.peer_id);
      const response = await fetch(`http://localhost:8081/connect-proxy?val=${host.peer_id}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          hostName: host.name,
          hostLocation: host.location,
          hostPeerID: host.peer_id,
          timestamp: new Date().toLocaleString(),
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to notify backend about the connection');
      }

      console.log(`Successfully notified backend about the connection to ${host.location}`);
    } catch (error) {
      console.error('Error notifying backend:', error);
    }
  };

  const handleAddProxy = async () => {
    if (newProxy.location.trim() === '' || newProxy.price.trim() === '') {
      alert('Please fill in all fields.');
      return;
    }
    await sendData();

    setProxyHosts([...proxyHosts, { ...newProxy, logs: [], isEnabled: false }]);

    setNewProxy({
      name: '',
      location: '',
      logs: [],
      peer_id: '',
      isEnabled: false,
      price: '',
      isHost: false,
    });

    setShowForm(false);
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

  const handleClearAndShowHistory = () => {
    setShowHistoryOnly(true);
  };

  const handleReturn = () => {
    setShowHistoryOnly(false);
  };

  return (
    <div className={styles.container}>
      <Box className={styles.boxContainer}>
        <Sidebar />
        <Box sx={{ marginTop: 2 }}>
          <Typography variant="h4">Proxy</Typography>
          <Box className={styles.header}>
            <Typography variant="h6">Your Current IP: {currentIP}</Typography>
            {!showHistoryOnly && (
              <Button variant="outlined" onClick={handleClearAndShowHistory}>
                Show Proxy History
              </Button>
            )}
          </Box>
          <br />
          {showHistoryOnly ? (
            <>
              <Box className={styles.historyContainer}>
                <Typography variant="h5">Proxy Connection History</Typography>
                <TableContainer className={styles.historyTable}>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell>Location</TableCell>
                        <TableCell>Connected At</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {proxyHistory.map((entry, index) => (
                        <TableRow key={index}>
                          <TableCell>{entry.name}</TableCell>
                          <TableCell>{entry.location}</TableCell>
                          <TableCell>{entry.timestamp}</TableCell>
                          <TableCell>
                            <Button
                              variant="contained"
                              onClick={() => {
                                const host = proxyHosts.find(h => h.location === entry.location);
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
                <Button variant="contained" onClick={handleReturn}>Return to Proxy Hosts</Button>
              </Box>
            </>
          ) : (
            <>
              <Box className={styles.proxyButton}>
                <Button variant="contained" onClick={() => setShowForm(prev => !prev)}>
                  {showForm ? 'Hide Form' : 'Add Yourself as Proxy'}
                </Button>

                <Box sx={{ display: 'flex', gap: '10px' }}>
                  <Button variant="outlined" onClick={handleSortByLocation}>Sort by Location</Button>
                  <Button variant="outlined" onClick={handleSortByPrice}>Sort by Price</Button>
                </Box>
              </Box>

              {showForm && (
                <Box sx={{ marginTop: 1 }} className={styles.form}>
                  <Typography variant="h6">Fill in your proxy details</Typography>
                  <Box sx={{ display: 'flex', gap: 2 }}>
                    <TextField
                      label="Name"
                      variant="outlined"
                      value={newProxy.name}
                      onChange={(e) => setNewProxy({ ...newProxy, name: e.target.value })}
                    />
                    <TextField
                      label="Location"
                      variant="outlined"
                      value={newProxy.location}
                      onChange={(e) => setNewProxy({ ...newProxy, location: e.target.value })}
                    />
                    <TextField
                      label="Price ($)"
                      variant="outlined"
                      type="number"
                      value={newProxy.price}
                      onChange={(e) => setNewProxy({ ...newProxy, price: e.target.value })}
                      InputProps={{ inputProps: { min: 0 } }}
                    />
                  </Box>
                  <br />
                  <Button variant="contained" onClick={handleAddProxy}>Add Proxy</Button>
                </Box>
              )}

              <TableContainer className={styles.proxyTable}>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell>Location</TableCell>
                      <TableCell>Price</TableCell>
                      <TableCell>Logs</TableCell>
                      <TableCell>Action</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {proxyHosts.map((host, index) => (
                      <TableRow key={index}>
                        <TableCell>{host.name}</TableCell>
                        <TableCell>{host.location}</TableCell>
                        <TableCell>{host.price}</TableCell>
                        <TableCell>
                          <Button variant="text" onClick={() => alert(host.logs.join('\n'))}>
                            View Logs
                          </Button>
                        </TableCell>
                        <TableCell>
                          <Button variant="contained" onClick={() => handleConnect(host)}>Connect</Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </>
          )}
        </Box>
      </Box>
    </div>
  );
};

export default ProxyHosts;
