Feature: Misc

    Scenario: Vocab
        When I use Vocab "Basic"
        Then the output should contain exactly:
            """
            test.md:3:5:Vale.Avoid:Avoid using 'Mac OS X'.
            test.md:7:1:Vale.Terms:Use 'definately' instead of 'Definately'.
            test.md:13:1:Vale.Terms:Use 'Documentarians' instead of 'documentarians'.
            test.md:17:11:Vale.Terms:Use 'Log4j' instead of 'log4j'.
            test.md:19:1:Vale.Terms:Use 'ABCDEF' instead of 'ABCDEf'.
            test.md:21:1:Vale.Terms:Use 'ABC-DEF' instead of 'ABC-DEf'.
            """

    Scenario: Multiple Vocabs
        When I use Vocab "Multi"
        Then the output should contain exactly:
            """
            test.md:3:5:Vale.Avoid:Avoid using 'Mac OS X'.
            test.md:7:29:Vale.Terms:Use 'IDEA' instead of 'idea'.
            test.md:13:1:Vale.Terms:Use 'Documentarians' instead of 'documentarians'.
            """

    Scenario: Line Endings
        When I test "misc/line-endings"
        Then the output should contain exactly:
            """
            CR.md:3:11:vale.Editorializing:Consider removing 'very'
            CR.md:32:11:vale.Editorializing:Consider removing 'very'
            CRLF.md:3:11:vale.Editorializing:Consider removing 'very'
            CRLF.md:32:11:vale.Editorializing:Consider removing 'very'
            """

    Scenario: Unicode
        When I test "misc/unicode"
        Then the output should contain exactly:
            """
            test.py:2:27:vale.Hedging:Consider removing 'probably'
            test.rst:1:36:vale.Annotations:'TODO' left in text
            test.rst:2:55:vale.Hedging:Consider removing 'probably'
            test.rst:4:23:vale.Editorializing:Consider removing 'very'
            test.txt:1:36:vale.Annotations:'TODO' left in text
            test.txt:2:55:vale.Hedging:Consider removing 'probably'
            test.txt:4:23:vale.Editorializing:Consider removing 'very'
            """

    Scenario: Duplicate matches
        When I test "misc/duplicates"
        Then the output should contain exactly:
            """
            test.adoc:1:3:Vale.Spelling:Did you really mean 'hav'?
            test.adoc:1:21:write-good.E-Prime:Avoid using "that's"
            test.adoc:4:3:Vale.Spelling:Did you really mean 'hav'?
            test.adoc:4:21:write-good.E-Prime:Avoid using "that's"
            test.md:1:6:write-good.E-Prime:Avoid using "is"
            test.md:1:9:vale.Editorializing:Consider removing 'very'
            test.md:1:15:vale.Editorializing:Consider removing 'very'
            test.md:1:39:vale.Editorializing:Consider removing 'very'
            test.md:1:57:vale.Editorializing:Consider removing 'very'
            test.md:3:14:write-good.E-Prime:Avoid using "are"
            test.md:3:28:write-good.E-Prime:Avoid using "be"
            test.md:5:9:write-good.E-Prime:Avoid using "is"
            test.md:5:95:Vale.Spelling:Did you really mean 'multimicrophone'?
            test.md:5:193:demo.CommasPerSentence:More than 3 commas!
            test2.md:1:6:write-good.E-Prime:Avoid using "is"
            test2.md:2:24:demo.Fancy:Change
            test2.md:2:40:write-good.E-Prime:Avoid using "there's"
            """

    Scenario: Nested markup
        When I test "misc/markup"
        Then the output should contain exactly:
            """
            test.md:4:1:Markup.Repetition:"in" is repeated.
            test.md:50:11:Markup.SentSpacing:"d.A" must contain one and only one space.
            """

    Scenario: Spelling
        When I test "spelling"
        Then the output should contain exactly:
            """
            test.adoc:61:1:Vale.Spelling:Did you really mean 'Nginx'?
            test.dic.md:1:3:Spelling.Test:'Test' is a typo!
            test.html:5:21:Spelling.Ignores:Did you really mean 'docbook'?
            test.html:14:96:Spelling.Ignores:Did you really mean 'TODO'?
            test.md:3:1:Spelling.Ignore:'HTTPie' is a typo!
            test.md:3:96:Spelling.Ignore:'human-friendly' is a typo!
            """

    Scenario: i18n
        When I test "i18n"
        Then the output should contain exactly:
            """
            zh.md:6:1:ZH.Simple:Avoid using "根据"
            zh.md:7:43:ZH.Simple:Avoid using "根据"
            """

    Scenario: infostrings
        When I test "misc/infostring"
        Then the output should contain exactly:
            """
            test.md:16:24:Vale.Spelling:Did you really mean 'Encryptor'?
            test.md:16:35:Vale.Spelling:Did you really mean 'Jasypt'?
            test.md:17:25:Vale.Spelling:Did you really mean 'config'?
            test.md:23:78:Vale.Spelling:Did you really mean 'config'?
            test.md:23:85:Vale.Spelling:Did you really mean 'json'?
            """
