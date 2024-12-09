import socket
import ssl
import threading
import logging
import select
import ipaddress
import argparse

# Configure logging
logging.basicConfig(level=logging.INFO, 
                    format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class HTTPSProxy:
    def __init__(self, host='127.0.0.1', port=8888, 
                 cert_file=None, key_file=None):
        """
        Initialize the HTTPS proxy server.
        
        :param host: Host IP to bind the proxy server
        :param port: Port number to listen on
        :param cert_file: SSL certificate file (optional)
        :param key_file: SSL private key file (optional)
        """
        self.host = host
        self.port = port
        self.server_socket = None
        
        # SSL context for server-side encryption if certificates provided
        self.ssl_context = None
        if cert_file and key_file:
            self.ssl_context = ssl.create_default_context(ssl.Purpose.CLIENT_AUTH)
            self.ssl_context.load_cert_chain(certfile=cert_file, keyfile=key_file)

    def start(self):
        """
        Start the proxy server and listen for incoming connections.
        """
        try:
            # Create a socket object
            self.server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            
            # Bind and listen
            self.server_socket.bind((self.host, self.port))
            self.server_socket.listen(100)
            
            logger.info(f"HTTPS Proxy server started on {self.host}:{self.port}")
            
            while True:
                # Accept incoming client connection
                client_socket, client_address = self.server_socket.accept()
                logger.info(f"Connection from {client_address}")
                
                # Create a new thread to handle the client
                client_thread = threading.Thread(
                    target=self.handle_client, 
                    args=(client_socket,)
                )
                client_thread.daemon = True
                client_thread.start()
        
        except Exception as e:
            logger.error(f"Error starting proxy server: {e}")
        finally:
            if self.server_socket:
                self.server_socket.close()

    def handle_client(self, client_socket):
        """
        Handle client connection and proxy requests.
        
        :param client_socket: Socket connection from client
        """
        try:
            # Receive the initial request
            request = client_socket.recv(4096)
            
            # Parse the CONNECT request
            if not request.startswith(b'CONNECT'):
                logger.warning("Not a CONNECT request")
                client_socket.close()
                return
            
            # Extract destination host and port
            destination = request.split(b' ')[1].decode()
            host, port = destination.split(':')
            port = int(port)
            
            # Validate host and port
            try:
                ipaddress.ip_address(host)
            except ValueError:
                # If not a valid IP, assume it's a domain name
                pass
            
            logger.info(f"Tunneling to {host}:{port}")
            
            # Establish connection to destination server
            try:
                destination_socket = socket.create_connection((host, port))
            except Exception as e:
                logger.error(f"Could not connect to destination: {e}")
                client_socket.send(b'HTTP/1.1 503 Service Unavailable\r\n\r\n')
                client_socket.close()
                return
            
            # Send successful connection response to client
            client_socket.send(b'HTTP/1.1 200 Connection Established\r\n\r\n')
            
            # If SSL context is set, wrap client socket
            if self.ssl_context:
                try:
                    client_socket = self.ssl_context.wrap_socket(
                        client_socket, 
                        server_side=True
                    )
                except ssl.SSLError as e:
                    logger.error(f"SSL error: {e}")
                    client_socket.close()
                    destination_socket.close()
                    return
            
            # Tunnel bidirectional traffic
            self.tunnel(client_socket, destination_socket)
        
        except Exception as e:
            logger.error(f"Error handling client: {e}")
        finally:
            client_socket.close()

    def tunnel(self, client_socket, destination_socket):
        """
        Tunnel bidirectional traffic between client and destination.
        
        :param client_socket: Client-side socket
        :param destination_socket: Destination server socket
        """
        def forward(source, destination):
            try:
                while True:
                    readable, _, _ = select.select([source], [], [], 5)
                    if not readable:
                        break
                    
                    data = source.recv(4096)
                    if not data:
                        break
                    
                    destination.sendall(data)
            except Exception as e:
                logger.error(f"Tunneling error: {e}")
        
        # Create threads for bidirectional forwarding
        client_to_dest = threading.Thread(
            target=forward, 
            args=(client_socket, destination_socket)
        )
        dest_to_client = threading.Thread(
            target=forward, 
            args=(destination_socket, client_socket)
        )
        
        client_to_dest.daemon = True
        dest_to_client.daemon = True
        
        client_to_dest.start()
        dest_to_client.start()
        
        # Wait for threads to complete
        client_to_dest.join()
        dest_to_client.join()
        
        # Close sockets
        client_socket.close()
        destination_socket.close()

def main():
    """
    Create and start the HTTPS proxy server with command-line port selection.
    """
    # Set up argument parsing
    parser = argparse.ArgumentParser(description='HTTPS Proxy Server')
    parser.add_argument('-p', '--port', 
                        type=int, 
                        default=8888, 
                        help='Port number to listen on (default: 8888)')
    parser.add_argument('-H', '--host', 
                        default='0.0.0.0', 
                        help='Host IP to bind the proxy server (default: 127.0.0.1)')
    
    # Parse arguments
    args = parser.parse_args()
    
    # Create and start proxy
    proxy = HTTPSProxy(
        host=args.host, 
        port=args.port
    )
    proxy.start()

if __name__ == '__main__':
    main()
