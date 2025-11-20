/**
 * PROGRAM KONVERSI ANGKA ROMAWI KE DESIMAL
 *
 * Penjelasan solusi:
 * 1. Saya membuat peta simbol (`romanMap`) yang menyimpan nilai setiap huruf Romawi
 *    (I, V, X, L, C, D, M). Dengan peta ini, lookup nilai setiap karakter menjadi operasi O(1).
 * 2. Fungsi `romanToInt` membaca string dari kiri ke kanan. Untuk setiap simbol, fungsi
 *    juga melihat simbol di posisi setelahnya. Jika nilai simbol saat ini lebih kecil daripada
 *    simbol berikutnya, berarti kita menemui pola pengurangan (misal I sebelum V membentuk IV).
 *    Dalam kasus ini, kita menambahkan (next - current) ke total dan melompati satu indeks.
 *    Jika tidak, kita tinggal menambahkan nilai simbol saat ini ke total.
 * 3. Pendekatan ini dipilih karena linearis (O(n)) dan langsung mengikuti cara penulisan angka
 *    Romawi—cukup satu pass tanpa struktur data tambahan.
 * 4. Logika pengurangan bekerja karena aturan resmi hanya mengizinkan angka kecil tertentu
 *    berada di depan angka yang lebih besar (I sebelum V/X, X sebelum L/C, C sebelum D/M).
 *    Dengan membandingkan nilai simbol sekarang dan berikutnya, kita menangkap semua kasus ini
 *    tanpa harus mendaftarkan kombinasi secara manual.
 * 5. Bagian penting kode:
 *    - Validasi input memastikan hanya string non-kosong berisi karakter Romawi sah.
 *    - Loop utama menerapkan aturan tambah atau pengurangan.
 *    - Contoh pemakaian di bagian bawah memperlihatkan hasil sesuai contoh soal.
 */

const romanMap = new Map([
  ['I', 1],
  ['V', 5],
  ['X', 10],
  ['L', 50],
  ['C', 100],
  ['D', 500],
  ['M', 1000],
]);


function romanToInt(romanString) {
  if (typeof romanString !== 'string' || romanString.trim() === '') {
    throw new Error('Input harus berupa string Romawi dan tidak boleh kosong.');
  }

  let total = 0;
  const roman = romanString.toUpperCase();

  for (let i = 0; i < roman.length; i += 1) {
    const currentSymbol = roman[i];
    const currentValue = romanMap.get(currentSymbol);

    if (currentValue === undefined) {
      throw new Error(`Simbol tidak valid: ${currentSymbol}`);
    }

    const nextSymbol = roman[i + 1];
    const nextValue = romanMap.get(nextSymbol);

    if (nextValue !== undefined && currentValue < nextValue) {
      total += nextValue - currentValue;
      i += 1; 
    } else {
      total += currentValue;
    }
  }

  return total;
}

// Contoh penggunaan sesuai soal
const samples = ['III', 'LVIII', 'MCMXCIV', 'XIV', 'CDXLIV'];
for (const sample of samples) {
  console.log(`${sample} -> ${romanToInt(sample)}`);
}

module.exports = { romanToInt };

