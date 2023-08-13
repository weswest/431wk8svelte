Feature: Checks
    Scenario: Script
        When I test "checks/Script"
        Then the output should contain exactly:
            """
            test.md:29:1:Checks.ScriptRE:Consider inserting a new section heading at this point.
            test.md:39:1:Checks.ScriptRE:Consider inserting a new section heading at this point.
            """

    Scenario: Metric
        When I test "checks/Metric"
        Then the output should contain exactly:
            """
            test.md:1:1:Checks.MetricValue:This topic has 1.00 H2s in it.
            """

    Scenario: Conditional
        When I test "checks/Conditional"
        Then the output should contain exactly:
            """
            test.md:9:5:Checks.MultiCapture:'NFL' has no definition
            """

    Scenario: Occurrence
        When I test "checks/Occurrence"
        Then the output should contain exactly:
            """
            test.md:1:1:demo.ZeroOccurrence:No intro
            test.md:1:3:demo.MinCount:Content too short.
            test2.md:1:3:demo.MinCount:Content too short.
            """

    Scenario: SentenceCase
        When I test "checks/SentenceCase"
        Then the output should contain exactly:
            """
            test.html:6:1:demo.SentenceCase:'Something Weird is happening With Vale With Entities' should be sentence-cased
            test.md:5:3:demo.SentenceCase:'this isn't in sentence case' should be sentence-cased
            test.md:11:3:demo.SentenceCase:'This Does Not Comply' should be sentence-cased
            """

    Scenario: Repetition
        When I test "checks/Repetition"
        Then the output should contain exactly:
            """
            test.tex:31:21:Vale.Repetition:'not' is repeated!
            text.rst:6:17:Vale.Repetition:'as' is repeated!
            text.rst:15:7:Vale.Repetition:'and' is repeated!
            text.rst:16:22:Vale.Repetition:'on' is repeated!
            text.rst:20:13:Vale.Repetition:'be' is repeated!
            """

    Scenario: Capitalization
        When I test "checks/Capitalization"
        Then the output should contain exactly:
            """
            test2.md:1:3:demo.CustomCap:'Section about the MacBook pro ultimate edition' should be in title case
            """

    Scenario: Spelling
        When I test "checks/Spelling"
        Then the output should contain exactly:
            """
            test.md:16:1:Vale.Spelling:Did you really mean 'gitlab'?
            """

    Scenario: Existence
        When I test "checks/Existence"
        Then the output should contain exactly:
            """
            """

    Scenario: Substitution
        When I test "checks/Substitution"
        Then the output should contain exactly:
            """
            """

    Scenario: Sequence
        When I test "checks/Sequence"
        Then the output should contain exactly:
            """
            test.adoc:3:4:LanguageTool.WOULD_BE_JJ_VB:The infinitive 'write' after 'be' requries 'to'. Did you mean 'be great *to* write'?
            test.adoc:9:88:LanguageTool.WOULD_BE_JJ_VB:The infinitive 'come' after 'be' requries 'to'. Did you mean 'be available *to* come'?
            test.adoc:11:32:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.adoc:13:5:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.adoc:15:24:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.adoc:17:42:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.adoc:21:5:LanguageTool.APOS_ARE:Did you mean "endpoints" instead of "endpoint's"?
            test.md:3:4:LanguageTool.WOULD_BE_JJ_VB:The infinitive 'write' after 'be' requries 'to'. Did you mean 'be great *to* write'?
            test.md:9:88:LanguageTool.WOULD_BE_JJ_VB:The infinitive 'come' after 'be' requries 'to'. Did you mean 'be available *to* come'?
            test.md:11:32:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.md:13:5:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.md:15:24:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.md:17:42:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.md:21:5:LanguageTool.APOS_ARE:Did you mean "endpoints" instead of "endpoint's"?
            test.md:25:1:LanguageTool.Metadata:Use data and metadata as plural nouns.
            test.md:29:1:LanguageTool.Metadata:Use data and metadata as plural nouns.
            test.txt:3:4:LanguageTool.WOULD_BE_JJ_VB:The infinitive 'write' after 'be' requries 'to'. Did you mean 'be great *to* write'?
            test.txt:9:88:LanguageTool.WOULD_BE_JJ_VB:The infinitive 'come' after 'be' requries 'to'. Did you mean 'be available *to* come'?
            test.txt:11:32:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.txt:13:5:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.txt:15:24:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.txt:17:42:LanguageTool.OF_ALL_TIMES:In this context, the idiom needs to be spelled 'of all time'.
            test.txt:21:5:LanguageTool.APOS_ARE:Did you mean "endpoints" instead of "endpoint's"?
            test.txt:25:1:LanguageTool.Metadata:Use data and metadata as plural nouns.
            test.txt:29:1:LanguageTool.Metadata:Use data and metadata as plural nouns.
            """
