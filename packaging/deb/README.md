# Debian



For debian packages you will need to add the following layouts during the build



eirene/
   DEBIAN/control
   DEBIAN/postinst
   usr/local/bin/eirene
   lib/systemd/system/eirene.service

This will be wrapped during the build package process building


Note this is still a work in progress:

TODO: removal/purge on removal using dpkg
      cleanup of control files to list what we want
      copyright inclusuion

