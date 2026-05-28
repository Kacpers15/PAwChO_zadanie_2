# Sprawozdanie

## 1. Opis konfiguracji poszczególnych etapów potoku (Pipeline)
Opracowany łańcuch GitHub Actions automatyzuje proces budowy, testowania bezpieczeństwa oraz publikacji aplikacji pogodowej. Potok składa się z następujących kroków:

* **1. Przygotowanie środowiska:**
  * **1.1** Pobranie kodu źródłowego (`actions/checkout`)
  * **1.2** Konfiguracja emulatora architektur `QEMU`
  * **1.3** Instalacja narzędzia `Docker Buildx`
       
* **2. Autoryzacja w rejestrach:**
  * **2.1** Logowanie do Docker Hub przy użyciu `Repository Secrets`
  * **2.2** Logowanie do GitHub Container Registry (GHCR) przy użyciu `${{ secrets.GITHUB_TOKEN }}`
       
* **4. Analiza bezpieczeństwa (Skanowanie CVE):**
  * **4.1** Wykorzystany skaner pozwala na blokowanie niezabezpiecznoego obrazu przed publikacją (musi byc bez błędów `HIGH` oraz `CRITICAL`
  * **4.2** Obraz jest tymczasowo budowany w architekutrze `linus/amd64` w celu przeskanowania go
       
* **5. Budowa docelowa i wysyłka (Push):**
  * **5.1** Po pomyślnym teście CVE, uruchamiany jest właściwy proces budowy wieloarchitekturalnej dla platform `linux/amd64` oraz `linux/arm64`.

---

## 2. Strategia tagowania obrazów i danych cache
Obrazy publikowane w rejestrze GitHub otrzymują dwa rodzaje tagów:
* **2.1** Tag gałęzi `main`: Wskazuje na najnowszą, aktualną wersję aplikacji.
* **2.2** Krótki identyfikator commitu Git np. `sha-7e581c1`: Zapewnia niezmienność obrazu, raz zbudowany kontener przypisany do konkretnego SHA nigdy nie ulegnie nadpisaniu.

Dane cache są przechowywane w dedykowanym repozytorium na Docker Hubie pod tagiem powiązanym z gałęzią (`weather-app-cache:main`) z flagą `mode=max`.
* **2.3.** Przeniesienie cache do zewnętrznego rejestru zapobiega zaśmiecaniu głównego rejestru (GHCR) tymczasowymi warstwami budowy.
* **2.3.** `mode=max` w przeciwieństwie do trybu domyślnego `mode=min`, zapisuje warstwy pamięci podręcznej dla wszystkich etapów pośrednich. Dzięki temu przy kolejnych uruchomieniach potoku, pobrane zależności Go są błyskawicznie odtwarzane z cache, co znacząco skraca czas budowania.

---

## 3. Potwierdzenie poprawności działania
Łańcuch GitHub Actions został uruchomiony z powodzeniem. Test CVE za pomocą Trivy nie wykazał krytycznych zagrożeń, co pozwoliło na pełną publikację pakietu.

---

## 4. Polecenie do uruchomienia lokalnego:
```bash
docker run -d -p 8080:8080 ghcr.io/kacpers15/pawcho_zadanie_2:main
