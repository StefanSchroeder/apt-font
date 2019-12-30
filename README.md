# apt-font

apt-font is a utility that helps you install fonts on Debian-based
Linux distributions. 

There are two primary use-cases:

* You know the name of a font, but don't know the name of the package
that contains it.

apt-font -f "Josef"

will install the package 'fonts-ldco' that contains the font *Josef* 
in ~/.fonts.

* You have a Libreoffice document. You open it and it looks ugly, because
the used fonts are not installed.

apt-font -i ~/Documents/my_fancy_document.odt

will install the fonts that are used in the document 'my_fancy_document'. Of course under the assumption that the respective fonts are available in the debian repository.

When you provide a directory, the tool will traverse the directory tree and install the fonts used in all the Libreoffice Writer (ODT), Libreoffice Spreadsheet (ODS) and Libreoffice Presentation (ODG) files that are found.

You can use *-h* for the list of options.

With *-s* you print the list of all known fonts to *apt-file* and the
packages that contain the font.

# Caveat

Currently the installation is only local to your user; no root prvileges are
required.

If a font is requested, that is already available as a system-font it
is installed in .fonts nevertheless.

If a font is requested, that is already installed in .fonts, it is re-installed.

# Design

There is a JSON configuration file that lists the Debian packages for every font family that is available in Debian (fonts.json). *apt-font* will lookup the 
requested font there and download and extract the debian-package to a 
temporary folder. Then it will copy the *.ttf fonts from the temp folder
to ~/.fonts. That's it.

If you want to find the fonts from a document, the file is inspected for
the used fonts (using the library odtfindfont) and these are processed as
described above.

*apt-font* was written by Stefan Schroeder, because he likes to install Linux
distros frequently, but needs to keep the system in a workable state, so that 
all the documents always look nice and tidy.

# Roadmap

I would be kinda neat to also install fonts that are not in the Debian 
repos, but that's somewhat hard. Also, we should support system-wide install,
which is pretty straightforward with the current architecture.

License: See LICENSE file. Don't worry, it's MIT.




