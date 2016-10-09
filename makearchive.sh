#!/bin/sh

name=fwego
current_commit=`git rev-parse HEAD`
short_commit=`echo ${current_commit:0:7}`
archive_name="$name-$short_commit.tar.gz"
folder_name="$name-$current_commit"

pushd ..
cp -R $name $folder_name
rm -rf $folder_name/.git
tar cfzv $archive_name $folder_name
rm -rf $folder_name
popd
