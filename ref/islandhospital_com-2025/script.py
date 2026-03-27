from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.serialization import pkcs12

# Function to extract the certificate chain from .pfx
def extract_cert_chain_from_pfx(pfx_file, pfx_password):
    # Open the .pfx file
    with open(pfx_file, 'rb') as pfx_file_obj:
        pfx_data = pfx_file_obj.read()

    # Load the pfx data
    private_key, certificate, cert_chain = pkcs12.load_key_and_certificates(pfx_data, pfx_password.encode(), backend=default_backend())

    return private_key, certificate, cert_chain

# Function to write the private key and certificate chain to files
def write_cert_files(private_key, certificate, cert_chain, cert_file, key_file):
    # Write the private key to a file
    with open(key_file, 'wb') as key_file_obj:
        key_file_obj.write(private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.TraditionalOpenSSL,
            encryption_algorithm=serialization.NoEncryption()
        ))

    # Write the server certificate (the main certificate)
    with open(cert_file, 'wb') as cert_file_obj:
        cert_file_obj.write(certificate.public_bytes(serialization.Encoding.PEM))

    # Write the intermediate certificates (if they exist) to the same file, appending
    with open(cert_file, 'ab') as cert_file_obj:
        for cert in cert_chain:
            cert_file_obj.write(cert.public_bytes(serialization.Encoding.PEM))

    print(f"Private key and certificate chain saved to: {key_file} and {cert_file}")

# Main function to handle the process
def main():
    # Specify paths for the .pfx file, password, and output files
    pfx_file = r'C:\Projects\nest-vesalius-m\ref\islandhospital_com-2025\wildcard.islandhospital.com_2025.pfx'  # Adjust to your actual pfx file location
    pfx_password = 'islandhospital.com2025'  # Adjust to your PFX password
    cert_file = 'imcgo.islandhospital.com.crt'  # Output file for the combined certificate chain
    key_file = 'imcgo.islandhospital.com.key'  # Output file for the private key

    private_key, certificate, cert_chain = extract_cert_chain_from_pfx(pfx_file, pfx_password)
    write_cert_files(private_key, certificate, cert_chain, cert_file, key_file)

if __name__ == "__main__":
    main()