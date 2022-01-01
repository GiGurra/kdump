use aes_gcm::{Aes256Gcm, Nonce};
use aes_gcm::aes::cipher::consts::U12;
use aes_gcm::aead::{Aead, NewAead};

use rand_chacha;
use rand::prelude::*;
use rand_chacha::ChaCha20Rng;

// https://docs.rs/aes-gcm/latest/aes_gcm

pub struct Encrypted {
    pub nonce_hex_string: String,
    pub encrypted_hex_string: String,
}

#[derive(Debug, PartialEq, Clone)]
pub enum EncryptError {
    InvalidKey,
    CipherEncryptFailure,
}

#[derive(Debug, PartialEq, Clone)]
pub enum DecryptError {
    _InvalidNonceHexInput,
    _InvalidDataHexInput,
    _InvalidKey,
    _CipherDecryptFailure,
    _InvalidUtf8Data,
}

pub fn encrypt(input: &str, key: &[u8]) -> Result<Encrypted, EncryptError> {
    let mut rng: ChaCha20Rng = ChaCha20Rng::from_entropy();

    let nonce_bytes_as_u32: Vec<u32> = (0..3).map(|_| rng.next_u32()).collect();
    let nonce_bytes: Vec<u8> = nonce_bytes_as_u32.iter().flat_map(|x| x.to_be_bytes()).collect();

    let nonce = Nonce::<U12>::from_slice(&nonce_bytes);
    let cipher = Aes256Gcm::new_from_slice(key).map_err(|_| EncryptError::InvalidKey)?;
    let encrypted_bytes = cipher.encrypt(nonce, input.as_bytes()).map_err(|_| EncryptError::CipherEncryptFailure)?;

    Ok(Encrypted {
        nonce_hex_string: hex::encode(nonce.to_vec()),
        encrypted_hex_string: hex::encode(encrypted_bytes),
    })
}

pub fn _decrypt(input: &Encrypted, key: &[u8]) -> Result<String, DecryptError> {
    let nonce_bytes = hex::decode(&input.nonce_hex_string).map_err(|_| DecryptError::_InvalidNonceHexInput)?;
    let encrypted_bytes = hex::decode(&input.encrypted_hex_string).map_err(|_| DecryptError::_InvalidDataHexInput)?;
    let nonce = Nonce::<U12>::from_slice(nonce_bytes.as_ref());
    let cipher = Aes256Gcm::new_from_slice(key).map_err(|_| DecryptError::_InvalidKey)?;
    let decrypted_bytes = cipher.decrypt(nonce, encrypted_bytes.as_ref()).map_err(|_|DecryptError::_CipherDecryptFailure)?;

    String::from_utf8(decrypted_bytes).map_err(|_| DecryptError::_InvalidUtf8Data)
}
