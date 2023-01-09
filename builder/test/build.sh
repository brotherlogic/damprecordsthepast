#!/usr/bin/zsh
for i in {1..16}
do
    echo "Downloading $i"
    sleep 5
    curl -s "https://api.discogs.com/artists/2228/releases?page=$i" -o "https___api.discogs.com_artists_2228_releases_page=$i"
done