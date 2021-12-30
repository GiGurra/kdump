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

pub fn encrypt(input: &str, key: &[u8]) -> Encrypted {

    let mut rng = ChaCha20Rng::from_entropy();

    let nonce_bytes_as_u32: Vec<u32> = (0..3).map(|_| rng.next_u32()).collect();
    let nonce_bytes: Vec<u8> = nonce_bytes_as_u32.iter().flat_map(|x| x.to_be_bytes()).collect();

    let nonce = Nonce::<U12>::from_slice(&nonce_bytes);
    let cipher = Aes256Gcm::new_from_slice(key).unwrap();
    let encrypted_bytes = cipher.encrypt(nonce, input.as_bytes()).unwrap();

    return Encrypted {
        nonce_hex_string: hex::encode(nonce.to_vec()),
        encrypted_hex_string: hex::encode(encrypted_bytes),
    };
}

pub fn _decrypt(input: &Encrypted, key: &[u8]) -> String {

    let nonce_bytes = hex::decode(&input.nonce_hex_string).unwrap();
    let nonce = Nonce::<U12>::from_slice(nonce_bytes.as_ref());
    let cipher = Aes256Gcm::new_from_slice(key).unwrap();
    let decrypted_bytes = cipher.decrypt(nonce, hex::decode(&input.encrypted_hex_string).unwrap().as_ref()).unwrap();

    return String::from_utf8(decrypted_bytes).unwrap();
}
