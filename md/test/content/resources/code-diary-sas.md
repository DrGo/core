+++
date = "2017-05-02T14:28:26-05:00"
title = "Documentation parsing in SAS"
+++

Code Diary is an automatic documentation parser for SAS.
It allows users to easily write, maintain, collate and share source code documentation.
The Code Diary syntax consists of comments, which add a “*” to the standard SAS comment syntax as well as section flags that identify in which output section a specific comment should be placed.
The tool generates two output documents.
One output document includes the location of comments (source file and line number) to aid programmers.
The other document, which excludes these technical details, is meant for sharing with the rest of the research team.
The SAS macro is available on [our GitHub page](https://github.com/VaccineAndDrugEvaluationCentre/code-diary-sas) along with several examples and is published in [SoftwareX](https://www.sciencedirect.com/science/article/pii/S2352711018300669).
