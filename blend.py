
from itertools import count
import numpy as np
import random

def calcBlends1(UA, UB):
    running_sim = 0
    running_dissim = 0
    UA_total = sum(UA.values())
    UB_total = sum(UB.values())

    percentARunning = 0
    percentBRunning = 0
    intersectingKeys = {}
    for artist_i in UA.keys():
        try:
            if UB.get(artist_i) != None:
                intersectingKeys[artist_i] = UA.get(artist_i) + UB.get(artist_i)
                print("", artist_i, ": ", UA.get(artist_i), " | ", UB.get(artist_i),  " | ", UB.get(artist_i) + UA.get(artist_i))

                percentA = UA.get(artist_i) / UA_total
                percentARunning += percentA
                percentB = UB.get(artist_i) / UB_total
                percentBRunning += percentB
        
        except:
            None

    sumIntersecting = sum(intersectingKeys.values())

    print("UA Total time", UA_total)
    print("UB Total time", UB_total)
    print("Total Intersecting time", sumIntersecting)
    print("Total Non-Intersecting time", UA_total + UB_total - sumIntersecting)



    print("Intersecting percent additive blend:", ((percentARunning + percentBRunning)))
    print("Intersecting Blend:", (sumIntersecting /( UA_total + UB_total)))

    print("Random multiplication of the above two:", ((percentARunning + percentBRunning) * sumIntersecting /( UA_total + UB_total)))


def calcBlends2(UA, UB):
    running_sim = 0
    running_dissim = 0
    UA_total = sum(UA.values())
    UB_total = sum(UB.values())

    percentARunning = 0
    percentBRunning = 0
    intersectingKeys = {}
    for artist_i in UA.keys():
        try:
            if UB.get(artist_i) != None:
                intersectingKeys[artist_i] = UA.get(artist_i) + UB.get(artist_i)
        except:
            None

    nonIntersectingKeys = UA | UB

    sumIntersecting = sum(intersectingKeys.values())
    sumNonIntersect = sum(nonIntersectingKeys.values())

    countIntersect = len(intersectingKeys)
    countUnion = len(nonIntersectingKeys)

    

    print(sumIntersecting/sumNonIntersect)  #Artist total time

    print(countIntersect, " / ", countUnion)    # Artist Int over Union
    print(countIntersect/countUnion) 

    frac_num = np.log( (countIntersect  )/countUnion) 

def self1(UA, UB):
    common_keys = set(UA.keys() & set(UB.keys()))  

    summation_of_common = sum(UA[k] * UA[k] for k in common_keys)

    magnitude_A = np.sqrt(sum(v ** 2 for v in UA.values()))
    magnitude_B = np.sqrt(sum(v ** 2 for v in UB.values()))

    if magnitude_A + magnitude_B == 0:
        print(0)
        return 0
    
    denom = (magnitude_A * magnitude_B)
    
    print(summation_of_common/denom)


def cosineSimilarity(UA, UB):
    common_keys = set(UA.keys()) & set(UB.keys())
    
    dot_product = sum(UA[k] * UB[k] for k in common_keys)
    
    magnitude_A = np.sqrt(sum(v ** 2 for v in UA.values()))
    magnitude_B = np.sqrt(sum(v ** 2 for v in UB.values()))
    
    if magnitude_A == 0 or magnitude_B == 0:
        print(str(0))
        return
    
    print( dot_product / (magnitude_A * magnitude_B))
    return dot_product / (magnitude_A * magnitude_B)
            
def weightedLogJaccard(UA, UB):
    keys = set(UA.keys()) | set(UB.keys())
    numerator = (sum(min(UA.get(k, 0), UB.get(k, 0)) for k in keys)) ** 1.0
    denominator = (sum(max(UA.get(k, 0), UB.get(k, 0)) for k in keys)) ** 1.0
    weighted_jaccard = numerator / denominator
    print(weighted_jaccard)
    return weighted_jaccard

def idf_log_weighted_cosine_similarity(UA, UB): 
    LAMBDA = 0.2
    common_keys = set(UA.keys()) & set(UB.keys())
    
    dot_product = 0
    for k in common_keys:
        termA = (UA[k]) 
        termB = (UB[k])
        dot_product += termA * termB


    magnitude_A = np.sqrt(sum(v ** 2 for v in UA.values()))
    magnitude_B = np.sqrt(sum(v ** 2 for v in UB.values()))
    
    if magnitude_A == 0 or magnitude_B == 0:
        print(str(0))
        return
    

    num_factor = np.log10(dot_product)
    denom_factor = np.log10(magnitude_A * magnitude_B)

    print("Log Weighted Value", num_factor/denom_factor)
    print("Direct Cosine Value", dot_product / (magnitude_A * magnitude_B))

    #Self 
    logWeightedValue = (num_factor/denom_factor)
    logWeightedValue = logWeightedValue if logWeightedValue > 0 else 0
    directCosineValue = dot_product/(magnitude_A*magnitude_B)
    finalValue = LAMBDA * logWeightedValue + (1 - LAMBDA) * directCosineValue

    print("Log Weighted Value", logWeightedValue)
    print("Direct Cosine Value", directCosineValue)
    print("Lambda Weighted Combination Value: ", finalValue)


def printStats(UA, UB): 
    print("Num UA: ", len(UA))
    print("Num UB: ", len(UB))

    common_keys = set(UA.keys()) & set(UB.keys())
    for k in common_keys: print(f'{k} - {UA.get(k)} - {UB.get(k)} = {UA.get(k) + UB.get(k)}')
    print("Sum of Common Keys", sum(UA.get(k) + UB.get(k) for k in common_keys))
    UA_total = sum(UA.values())
    UB_total = sum(UB.values())
    print("UA Total: ", UA_total)
    print("UB Total: ", UB_total)
    print('Common Time: ', sum(UA[k] + UB[k] for k in common_keys))
    print("Total: ", UA_total + UB_total)



#Bunch of random
UA = {'K': 8109, 'P': 1837, 'S': 4670, 'C': 7593, 'L': 7381, 'Y': 1254, 'E': 2386, 'I': 1651, 'Z': 9967, 'U': 3478}
UB = {'G': 6614, 'V': 7075, 'U': 1634, 'W': 306, 'E': 908, 'I': 6948, 'D': 1903, 'A': 4431, 'K': 1919, 'R': 2158}
UB_low  = {'G': 6614, 'V': 7075, 'U': 1634, 'W': 306, 'D': 1903, 'A': 4431, 'K': 1919, 'R': 2158}
UA_high  = {'J': 938, 'P': 1820, 'D': 6439, 'A': 9953, 'V': 6568, 'W': 1835, 'Z': 8655, 'N': 9600, 'E': 366, 'U': 6808, 'G': 1977, 'R': 9837, 'S': 1583, 'I': 6298, 'M': 3559, 'L': 7330, 'T': 3058, 'Q': 5624, 'O': 4467, 'F': 8084, 'Y': 4406, 'C': 8714, 'B': 9073}
UB_single = {'G': 1000}
UB_time_low = {'G': 664, 'V': 705, 'U': 164, 'W': 106, 'E': 108, 'I': 1948, 'D': 903, 'A': 431, 'K': 919, 'R': 158}
UB_double = {'J': 938, 'P': 1820}
UB_time_very_low = {'G': 4, 'V': 70, 'U': 14, 'W': 10, 'E': 18, 'I': 194, 'D': 90, 'A': 41, 'K': 99, 'R': 18}




def test(UA, UB, function_to_test):
    printStats(UA, UB)
    function_to_test(UA, UB)

def test_double(UA, UB, f1, f2):
    printStats(UA, UB)
    first = f1(UA, UB)
    second = f2(UA, UB)
    print((first + second)/2)

test(UB_double, UB_single, idf_log_weighted_cosine_similarity )

# test_double(UA, UB_low, cosineSimilarity, weightedLogJaccard)




